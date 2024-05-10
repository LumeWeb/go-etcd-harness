package etcd_harness_test

import (
	"os"
	"testing"
	"time"

	"github.com/chen-anders/go-etcd-harness"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	etcd "go.etcd.io/etcd/client/v3"
	"golang.org/x/net/context"
)

type HarnessTestSuite struct {
	suite.Suite
	kv etcd.KV
}

func (s *HarnessTestSuite) SetupSuite() {
	_, err := s.kv.Put(newContext(), "/testdir2", "", etcd.WithPrevKV())
	require.NoError(s.T(), err, "creating the test directory must never fail.")
}

func (s *HarnessTestSuite) TestReadWrite() {
	_, err := s.kv.Put(newContext(), "/testdir/somevalue", "SomeContent")
	require.NoError(s.T(), err, "set must succeed")
	resp, err := s.kv.Get(newContext(), "/testdir/somevalue")
	require.NoError(s.T(), err, "get must succeed")
	assert.Equal(s.T(), "SomeContent", string(resp.Kvs[0].Value))
}

func TestHarnessTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skipf("HarnessTestSuite is a long integration test suite. Skipping due to test short.")
	}
	if !etcd_harness.LocalEtcdAvailable() {
		t.Skipf("etcd is not available in $PATH, skipping suite")
	}

	harness, err := etcd_harness.New(os.Stderr)
	if err != nil {
		t.Fatalf("failed starting etcd harness: %v", err)
	}
	t.Logf("will use etcd harness endpoint: %v", harness.Endpoint)
	defer func() {
		harness.Stop()
		t.Logf("cleaned up etcd harness")
	}()
	suite.Run(t, &HarnessTestSuite{kv: etcd.NewKV(harness.Client)})
}

func newContext() context.Context {
	c, _ := context.WithTimeout(context.TODO(), 1*time.Second)
	return c
}
