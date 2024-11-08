package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommand(t *testing.T) {
	cmdReq := CreateCommandRequestParams{
		InputContent: "nextjs init",
		CommandType:  "question",
	}

	cmd, err := testStore.CreateCommandRequest(context.Background(), cmdReq)
	fmt.Println(cmd)

	require.NoError(t, err)
	require.NotEmpty(t, cmd)
	require.Equal(t, cmdReq.InputContent, cmd.InputContent)
}
