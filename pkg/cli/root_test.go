package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/obot-platform/nanobot/pkg/config"
	"github.com/obot-platform/nanobot/pkg/types"
)

func TestNanobotConfigPathsDefault(t *testing.T) {
	n := &Nanobot{}

	paths := n.ConfigPaths()
	if len(paths) != 1 || paths[0] != config.DefaultConfigPath {
		t.Fatalf("expected default config path [.nanobot/], got %v", paths)
	}
}

func TestRuntimeConfigDirDefaultsWhenNoPathsExist(t *testing.T) {
	configDir := runtimeConfigDir([]string{"./missing", "https://example.com/nanobot.yaml"})
	if configDir != config.DefaultConfigPath {
		t.Fatalf("expected default config dir %q, got %q", config.DefaultConfigPath, configDir)
	}
}

func TestRuntimeConfigDirKeepsDefaultForExistingDirectory(t *testing.T) {
	dir := t.TempDir()

	configDir := runtimeConfigDir([]string{"./missing", dir})
	if configDir != config.DefaultConfigPath {
		t.Fatalf("expected default config dir %q, got %q", config.DefaultConfigPath, configDir)
	}
}

func TestRuntimeConfigDirReturnsFileParentDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nanobot.yaml")
	if err := os.WriteFile(path, []byte("agents: {}\n"), 0o644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	configDir := runtimeConfigDir([]string{path})
	if configDir != dir {
		t.Fatalf("expected config dir %q, got %q", dir, configDir)
	}
}

func TestBuiltinAimlapiProvider(t *testing.T) {
	n := &Nanobot{}

	provider, ok := n.llmConfig().LLMProviders["aimlapi"]
	if !ok {
		t.Fatal("expected a built-in aimlapi provider")
	}

	if provider.Dialect != types.DialectOpenAIChatCompletions {
		t.Errorf("dialect: got %q, want %q", provider.Dialect, types.DialectOpenAIChatCompletions)
	}

	// An empty base URL is not a neutral default here: the chat-completions
	// client substitutes api.openai.com for it, which would send an AI/ML API
	// key to OpenAI. The built-in entry must therefore carry a literal.
	if provider.BaseURL != "https://api.aimlapi.com/v1" {
		t.Errorf("baseURL: got %q, want the AI/ML API endpoint", provider.BaseURL)
	}

	if provider.APIKey != "${AIMLAPI_API_KEY}" {
		t.Errorf("apiKey: got %q, want the AIMLAPI_API_KEY reference", provider.APIKey)
	}

	if got := provider.Headers["X-AIMLAPI-Source"]; got != "agent/nanobot" {
		t.Errorf("X-AIMLAPI-Source: got %q, want %q", got, "agent/nanobot")
	}
}
