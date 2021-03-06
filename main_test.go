package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
	"strings"
	"testing"
)

func getMockOsExiter(actualExitCodeChannel chan int) func(int) {
	return func(actualExitCode int) {
		actualExitCodeChannel <- actualExitCode
	}
}

func testMainWithArguments(t *testing.T, expectedExitCode int, args ...string) {
	actualExitCodeChannel := make(chan int)
	cli.OsExiter = getMockOsExiter(actualExitCodeChannel)

	appArgs := append([]string{"name-of-binary"}, args...)
	go RunApp(appArgs)
	actualExitCode := <-actualExitCodeChannel

	assert.Equal(t, expectedExitCode, actualExitCode, strings.Join(args, " "))
}

//func TestRunMain(t *testing.T) {
//testMainWithArguments(t, okExitCode)
//}

func TestProcessMailDefault(t *testing.T) {
	testMainWithArguments(t, okExitCode, "process-mail")
}

func TestProcessMailInvalidPassphrase(t *testing.T) {
	testMainWithArguments(t, errorExitCode, "process-mail", "--passphrase", "invalid!")
}

func TestProcessMailPrivateKeyInvalid(t *testing.T) {
	testMainWithArguments(t, errorExitCode, "process-mail", "--private-key", "./test/mails/plaintext.eml")
}

func TestProcessMailPrivateKeyNotFound(t *testing.T) {
	testMainWithArguments(t, errorExitCode, "process-mail", "--private-key", "does_not_exist")
}

func testProcessMail(t *testing.T, expectedExitCode int, file string, args ...string) {
	filePath := fmt.Sprintf("./test/mails/%s", file)
	mainArgs := append([]string{"process-mail", "--passphrase", "validation", "--file", filePath}, args...)
	testMainWithArguments(t, expectedExitCode, mainArgs...)
}

func TestProcessMailFilesSuccessfully(t *testing.T) {
	testProcessMail(t, okExitCode, "attachment.eml")
	testProcessMail(t, okExitCode, "crypted_signed_request_enigmail.eml")
	testProcessMail(t, okExitCode, "plaintext.eml")
	testProcessMail(t, okExitCode, "signed_multipart_simple.eml")
	testProcessMail(t, okExitCode, "signed_request_enigmail.eml")
}

func TestProcessFileError(t *testing.T) {
	testProcessMail(t, errorExitCode, "invalid")
}

func TestProcessMailStorageTypes(t *testing.T) {
	testProcessMail(t, okExitCode, "attachment.eml", "--storage", "none")
	testProcessMail(t, okExitCode, "attachment.eml", "--storage", "memory")
	testProcessMail(t, okExitCode, "attachment.eml", "--storage", "file")

	testProcessMail(t, errorExitCode, "attachment.eml", "--storage", "invalid")
}
