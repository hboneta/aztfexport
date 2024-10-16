package resource

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Azure/aztfexport/internal/client"
	internalconfig "github.com/Azure/aztfexport/internal/config"

	"github.com/Azure/aztfexport/pkg/config"

	"github.com/Azure/aztfexport/internal"
	"github.com/Azure/aztfexport/internal/test"
	"github.com/Azure/aztfexport/internal/test/cases"
	"github.com/Azure/aztfexport/internal/utils"
	"github.com/hashicorp/terraform-exec/tfexec"
)

func runCase(t *testing.T, d test.Data, c cases.Case) {
	tfexecPath := test.EnsureTF(t)

	provisionDir := t.TempDir()
	if test.Keep() {
		provisionDir, _ = os.MkdirTemp("", "")
		t.Log(provisionDir)
	}

	if err := os.WriteFile(filepath.Join(provisionDir, "main.tf"), []byte(c.Tpl(d)), 0644); err != nil {
		t.Fatalf("created to create the TF config file: %v", err)
	}
	tf, err := tfexec.NewTerraform(provisionDir, tfexecPath)
	if err != nil {
		t.Fatalf("failed to new terraform: %v", err)
	}
	ctx := context.Background()
	t.Log("Running: terraform init")
	if err := tf.Init(ctx); err != nil {
		t.Fatalf("terraform init failed: %v", err)
	}
	t.Log("Running: terraform apply")
	if err := tf.Apply(ctx); err != nil {
		t.Fatalf("terraform apply failed: %v", err)
	}

	if !test.Keep() {
		defer func() {
			t.Log("Running: terraform destroy")
			if err := tf.Destroy(ctx); err != nil {
				t.Logf("terraform destroy failed: %v", err)
			}
		}()
	}

	const delay = time.Minute
	t.Logf("Sleep for %v to wait for the just created resources be recorded in ARG\n", delay)
	time.Sleep(delay)

	aztfexportDir := t.TempDir()
	l, err := c.SingleResourceContext(d)
	if err != nil {
		t.Fatalf("failed to get resource ids: %v", err)
	}

	cred, clientOpt := test.BuildCredAndClientOpt(t)

	cfg := internalconfig.NonInteractiveModeConfig{
		Config: config.Config{
			CommonConfig: config.CommonConfig{
				Logger:               slog.New(slog.NewTextHandler(io.Discard, nil)),
				SubscriptionId:       os.Getenv("ARM_SUBSCRIPTION_ID"),
				AzureSDKCredential:   cred,
				AzureSDKClientOption: *clientOpt,
				OutputDir:            aztfexportDir,
				BackendType:          "local",
				DevProvider:          true,
				Parallelism:          1,
				ProviderName:         "azurerm",
			},
		},
		PlainUI: true,
	}

	// Test single resouce mode
	for idx, rctx := range l {
		cfg.ResourceIds = []string{rctx.AzureId}
		cfg.TFResourceName = fmt.Sprintf("res-%d", idx)

		t.Logf("Single resource importing %s\n", rctx.AzureId)
		if err := utils.RemoveEverythingUnder(cfg.OutputDir); err != nil {
			t.Fatalf("failed to clean up the output directory: %v", err)
		}
		if err := internal.BatchImport(ctx, cfg); err != nil {
			t.Fatalf("failed to run resource import: %v", err)
		}
		test.Verify(t, ctx, aztfexportDir, tfexecPath, rctx.ExpectResourceCount)
	}

	// Test multi-resource mode
	var resourceIds []string
	var total int
	for _, rctx := range l {
		resourceIds = append(resourceIds, rctx.AzureId)
		total += rctx.ExpectResourceCount
	}

	cfg.ResourceNamePattern = "res-"
	cfg.ResourceIds = resourceIds
	cfg.CommonConfig.Parallelism = total

	t.Logf("Multiple resources importing %v\n", resourceIds)
	if err := utils.RemoveEverythingUnder(cfg.OutputDir); err != nil {
		t.Fatalf("failed to clean up the output directory: %v", err)
	}
	if err := internal.BatchImport(ctx, cfg); err != nil {
		t.Fatalf("failed to run resource import: %v", err)
	}
	test.Verify(t, ctx, aztfexportDir, tfexecPath, total)
}

func TestVnet(t *testing.T) {
	t.Parallel()
	test.Precheck(t)
	c, d := cases.CaseVnet{}, test.NewData()
	runCase(t, d, c)
}

func TestComputeVMDisk(t *testing.T) {
	t.Parallel()
	test.Precheck(t)
	c, d := cases.CaseComputeVMDisk{}, test.NewData()
	runCase(t, d, c)
}

func TestSignalRService(t *testing.T) {
	t.Parallel()
	test.Precheck(t)
	c, d := cases.CaseSignalRService{}, test.NewData()
	runCase(t, d, c)
}

func TestApplicationInsightWebTest(t *testing.T) {
	t.Parallel()
	test.Precheck(t)
	c, d := cases.CaseApplicationInsightWebTest{}, test.NewData()
	runCase(t, d, c)
}

func TestKeyVaultNestedItems(t *testing.T) {
	t.Parallel()
	test.Precheck(t)
	cred, opt := test.BuildCredAndClientOpt(t)
	c, d := cases.CaseKeyVaultNestedItems{B: client.ClientBuilder{Credential: cred, Opt: *opt}}, test.NewData()
	runCase(t, d, c)
}

func TestFunctionAppSlot(t *testing.T) {
	t.Parallel()
	test.Precheck(t)
	c, d := cases.CaseFunctionAppSlot{}, test.NewData()
	runCase(t, d, c)
}

func TestStorageFileShare(t *testing.T) {
	t.Parallel()
	test.Precheck(t)
	c, d := cases.CaseStorageFileShare{}, test.NewData()
	runCase(t, d, c)
}
