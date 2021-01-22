package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/tblyler/httpwrap/config"
)

// CommandResponse to be returned from an endpoint
type CommandResponse struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	StdOut    string    `json:"stdout,omitempty"`
	StdErr    string    `json:"stderr,omitempty"`
	ExitCode  int       `json:"exit_code,omitempty"`
	Error     string    `json:"error,omitempty"`
}

func main() {
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if configFilePath == "" {
		fmt.Fprintln(os.Stderr, "CONFIG_FILE_PATH environment variable must be set")
		os.Exit(1)
	}

	config, err := config.NewJSONFileSource(configFilePath).Config(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}

	for route, endpoint := range config.Endpoints {
		endpoint := endpoint
		http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
			if endpoint.HTTPMethod == "" {
				endpoint.HTTPMethod = http.MethodGet
			}

			if strings.ToUpper(endpoint.HTTPMethod) != strings.ToUpper(r.Method) {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "unconfigured HTTP method:", r.Method)
				return
			}

			args := endpoint.Arguments
			if endpoint.AllowExternalArguments {
				externalArgs, _ := r.URL.Query()["args"]
				args = append(args, externalArgs...)
			}

			command := exec.CommandContext(
				r.Context(),
				endpoint.Command,
				args...,
			)

			if endpoint.AllowStdin {
				command.Stdin = r.Body
			}

			stdErr := bytes.NewBuffer(nil)
			if !endpoint.DiscardStderr {
				command.Stderr = stdErr
			}

			stdOut := bytes.NewBuffer(nil)
			if !endpoint.DiscardStdout {
				command.Stdout = stdOut
			}

			response := &CommandResponse{
				StartTime: time.Now(),
			}

			err := command.Start()
			if err != nil {
				response.EndTime = time.Now()
				response.Error = fmt.Sprintf("failed to start command: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(response)
				return
			}

			command.Wait()

			response.EndTime = time.Now()
			response.ExitCode = command.ProcessState.ExitCode()
			response.StdErr = stdErr.String()
			response.StdOut = stdOut.String()

			json.NewEncoder(w).Encode(response)
		})
	}

	httpServer := &http.Server{
		Addr: fmt.Sprintf("%s:%d", config.ListenAddress, config.ListenPort),
	}

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	go func() {
		<-ctx.Done()
		fmt.Println("got interrupt signal, waiting up to 5 seconds to gracefully shut down")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(5))
		defer cancel()

		httpServer.Shutdown(ctx)
	}()

	fmt.Println("listening on:", httpServer.Addr)
	err = httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("got error from HTTP server", err.Error())
		os.Exit(3)
	}

	fmt.Println("successfully shut down")
}
