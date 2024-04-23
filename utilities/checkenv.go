package utilities

import "os"

func CheckEnv(envs []string) bool {
	msgs := []string{}

	for _, env := range envs {
		if _, ok := os.LookupEnv(env); !ok {
			msgs = append(msgs, env+" is required and not set!")
		}
	}

	for _, msg := range msgs {
		println(msg)
	}

	return len(msgs) == 0
}
