package commit

import (
	"github.com/jesseduffield/lazygit/pkg/config"
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

var fileModHook = `#!/bin/bash

if [[ -f test-wip-commit-prefix ]]; then
  echo "Modified text" > test-wip-commit-prefix
fi
`

var CommitSkipHooks = NewIntegrationTest(NewIntegrationTestArgs{
	Description:  "Commit with skip hook using CommitChangesWithoutHook",
	ExtraCmdArgs: []string{},
	Skip:         false,
	SetupConfig:  func(config *config.AppConfig) {},
	SetupRepo: func(shell *Shell) {
		shell.CreateFile(".git/hooks/pre-commit", fileModHook)
		shell.MakeExecutable(".git/hooks/pre-commit")

		shell.NewBranch("feature/TEST-002")
		shell.CreateFile("test-wip-commit-prefix", "Initial text")
	},
	Run: func(t *TestDriver, keys config.KeybindingConfig) {
		t.Views().Commits().
			IsEmpty()

		t.Views().Files().
			IsFocused().
			PressPrimaryAction().
			Press(keys.Files.CommitChangesWithoutHook)

		t.ExpectPopup().CommitMessagePanel().
			Title(Equals("Commit summary")).
			Type("foo bar").
			Confirm()

		t.FileSystem().FileContent("test-wip-commit-prefix", Equals("Initial text"))

		t.Views().Commits().Focus()
		t.Views().Main().Content(Contains("foo bar"))
		t.Views().Extras().Content(Contains("--no-verify"))
	},
})
