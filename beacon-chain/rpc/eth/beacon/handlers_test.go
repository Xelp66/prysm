package beacon

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	testing2 "github.com/prysmaticlabs/prysm/v4/beacon-chain/blockchain/testing"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/core/transition"
	dbTest "github.com/prysmaticlabs/prysm/v4/beacon-chain/db/testing"
	doublylinkedtree "github.com/prysmaticlabs/prysm/v4/beacon-chain/forkchoice/doubly-linked-tree"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/rpc/eth/shared"
	rpctesting "github.com/prysmaticlabs/prysm/v4/beacon-chain/rpc/eth/shared/testing"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/rpc/testutil"
	"github.com/prysmaticlabs/prysm/v4/beacon-chain/state"
	mockSync "github.com/prysmaticlabs/prysm/v4/beacon-chain/sync/initial-sync/testing"
	"github.com/prysmaticlabs/prysm/v4/config/params"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/blocks"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/interfaces"
	"github.com/prysmaticlabs/prysm/v4/encoding/bytesutil"
	"github.com/prysmaticlabs/prysm/v4/proto/migration"
	eth "github.com/prysmaticlabs/prysm/v4/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v4/testing/assert"
	mock2 "github.com/prysmaticlabs/prysm/v4/testing/mock"
	"github.com/prysmaticlabs/prysm/v4/testing/require"
	"github.com/prysmaticlabs/prysm/v4/testing/util"
	"github.com/stretchr/testify/mock"
)

func TestPublishBlockV2(t *testing.T) {
	ctrl := gomock.NewController(t)

	t.Run("Phase 0", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			block, ok := req.Block.(*eth.GenericSignedBeaconBlock_Phase0)
			converted, err := shared.BeaconBlockFromConsensus(block.Phase0.Block)
			require.NoError(t, err)
			var signedblock *shared.SignedBeaconBlock
			err = json.Unmarshal([]byte(rpctesting.Phase0Block), &signedblock)
			require.NoError(t, err)
			require.DeepEqual(t, converted, signedblock.Message)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(rpctesting.Phase0Block)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("Altair", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			block, ok := req.Block.(*eth.GenericSignedBeaconBlock_Altair)
			converted, err := shared.BeaconBlockAltairFromConsensus(block.Altair.Block)
			require.NoError(t, err)
			var signedblock *shared.SignedBeaconBlockAltair
			err = json.Unmarshal([]byte(rpctesting.AltairBlock), &signedblock)
			require.NoError(t, err)
			require.DeepEqual(t, converted, signedblock.Message)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(rpctesting.AltairBlock)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("Bellatrix", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			block, ok := req.Block.(*eth.GenericSignedBeaconBlock_Bellatrix)
			converted, err := shared.BeaconBlockBellatrixFromConsensus(block.Bellatrix.Block)
			require.NoError(t, err)
			var signedblock *shared.SignedBeaconBlockBellatrix
			err = json.Unmarshal([]byte(rpctesting.BellatrixBlock), &signedblock)
			require.NoError(t, err)
			require.DeepEqual(t, converted, signedblock.Message)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(rpctesting.BellatrixBlock)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("Capella", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			block, ok := req.Block.(*eth.GenericSignedBeaconBlock_Capella)
			converted, err := shared.BeaconBlockCapellaFromConsensus(block.Capella.Block)
			require.NoError(t, err)
			var signedblock *shared.SignedBeaconBlockCapella
			err = json.Unmarshal([]byte(rpctesting.CapellaBlock), &signedblock)
			require.NoError(t, err)
			require.DeepEqual(t, converted, signedblock.Message)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(rpctesting.CapellaBlock)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("Deneb", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			block, ok := req.Block.(*eth.GenericSignedBeaconBlock_Deneb)
			converted, err := shared.BeaconBlockDenebFromConsensus(block.Deneb.Block.Block)
			require.NoError(t, err)
			var signedblock *shared.SignedBeaconBlockContentsDeneb
			err = json.Unmarshal([]byte(rpctesting.DenebBlockContents), &signedblock)
			require.NoError(t, err)
			require.DeepEqual(t, converted, signedblock.SignedBlock.Message)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(rpctesting.DenebBlockContents)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("invalid block", func(t *testing.T) {
		server := &Server{
			SyncChecker: &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(rpctesting.BlindedBellatrixBlock)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlockV2(writer, request)
		assert.Equal(t, http.StatusBadRequest, writer.Code)
		assert.Equal(t, true, strings.Contains(writer.Body.String(), "Body does not represent a valid block type"))
	})
	t.Run("syncing", func(t *testing.T) {
		chainService := &testing2.ChainService{}
		server := &Server{
			SyncChecker:           &mockSync.Sync{IsSyncing: true},
			HeadFetcher:           chainService,
			TimeFetcher:           chainService,
			OptimisticModeFetcher: chainService,
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte("foo")))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlockV2(writer, request)
		assert.Equal(t, http.StatusServiceUnavailable, writer.Code)
		assert.Equal(t, true, strings.Contains(writer.Body.String(), "Beacon node is currently syncing and not serving request on that endpoint"))
	})
}

func TestPublishBlockV2SSZ(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Run("Bellatrix", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			_, ok := req.Block.(*eth.GenericSignedBeaconBlock_Bellatrix)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}
		var bellablock shared.SignedBeaconBlockBellatrix
		err := json.Unmarshal([]byte(rpctesting.BellatrixBlock), &bellablock)
		require.NoError(t, err)
		genericBlock, err := bellablock.ToGeneric()
		require.NoError(t, err)
		sszvalue, err := genericBlock.GetBellatrix().MarshalSSZ()
		require.NoError(t, err)
		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader(sszvalue))
		request.Header.Set("Accept", "application/octet-stream")
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("Capella", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			_, ok := req.Block.(*eth.GenericSignedBeaconBlock_Capella)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		var cblock shared.SignedBeaconBlockCapella
		err := json.Unmarshal([]byte(rpctesting.CapellaBlock), &cblock)
		require.NoError(t, err)
		genericBlock, err := cblock.ToGeneric()
		require.NoError(t, err)
		sszvalue, err := genericBlock.GetCapella().MarshalSSZ()
		require.NoError(t, err)
		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader(sszvalue))
		request.Header.Set("Accept", "application/octet-stream")
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("Deneb", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			_, ok := req.Block.(*eth.GenericSignedBeaconBlock_Deneb)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		var dblock shared.SignedBeaconBlockContentsDeneb
		err := json.Unmarshal([]byte(rpctesting.DenebBlockContents), &dblock)
		require.NoError(t, err)
		genericBlock, err := dblock.ToGeneric()
		require.NoError(t, err)
		v2block, err := migration.V1Alpha1SignedBeaconBlockDenebAndBlobsToV2(genericBlock.GetDeneb())
		require.NoError(t, err)
		sszvalue, err := v2block.MarshalSSZ()
		require.NoError(t, err)
		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader(sszvalue))
		request.Header.Set("Accept", "application/octet-stream")
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("invalid block", func(t *testing.T) {
		server := &Server{
			SyncChecker: &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(rpctesting.BlindedBellatrixBlock)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlockV2(writer, request)
		assert.Equal(t, http.StatusBadRequest, writer.Code)
		assert.Equal(t, true, strings.Contains(writer.Body.String(), "Body does not represent a valid block type"))
	})
}

func TestPublishBlindedBlockV2(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Run("Phase 0", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			block, ok := req.Block.(*eth.GenericSignedBeaconBlock_Phase0)
			converted, err := shared.BeaconBlockFromConsensus(block.Phase0.Block)
			require.NoError(t, err)
			var signedblock *shared.SignedBeaconBlock
			err = json.Unmarshal([]byte(rpctesting.Phase0Block), &signedblock)
			require.NoError(t, err)
			require.DeepEqual(t, converted, signedblock.Message)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(rpctesting.Phase0Block)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlindedBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("Altair", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			block, ok := req.Block.(*eth.GenericSignedBeaconBlock_Altair)
			converted, err := shared.BeaconBlockAltairFromConsensus(block.Altair.Block)
			require.NoError(t, err)
			var signedblock *shared.SignedBeaconBlockAltair
			err = json.Unmarshal([]byte(rpctesting.AltairBlock), &signedblock)
			require.NoError(t, err)
			require.DeepEqual(t, converted, signedblock.Message)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(rpctesting.AltairBlock)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlindedBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("Blinded Bellatrix", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			block, ok := req.Block.(*eth.GenericSignedBeaconBlock_BlindedBellatrix)
			converted, err := shared.BlindedBeaconBlockBellatrixFromConsensus(block.BlindedBellatrix.Block)
			require.NoError(t, err)
			var signedblock *shared.SignedBlindedBeaconBlockBellatrix
			err = json.Unmarshal([]byte(rpctesting.BlindedBellatrixBlock), &signedblock)
			require.NoError(t, err)
			require.DeepEqual(t, converted, signedblock.Message)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(rpctesting.BlindedBellatrixBlock)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlindedBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("Blinded Capella", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			block, ok := req.Block.(*eth.GenericSignedBeaconBlock_BlindedCapella)
			converted, err := shared.BlindedBeaconBlockCapellaFromConsensus(block.BlindedCapella.Block)
			require.NoError(t, err)
			var signedblock *shared.SignedBlindedBeaconBlockCapella
			err = json.Unmarshal([]byte(rpctesting.BlindedCapellaBlock), &signedblock)
			require.NoError(t, err)
			require.DeepEqual(t, converted, signedblock.Message)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(rpctesting.BlindedCapellaBlock)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlindedBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("Blinded Deneb", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			block, ok := req.Block.(*eth.GenericSignedBeaconBlock_BlindedDeneb)
			converted, err := shared.BlindedBeaconBlockDenebFromConsensus(block.BlindedDeneb.Block.Block)
			require.NoError(t, err)
			var signedblock *shared.SignedBlindedBeaconBlockContentsDeneb
			err = json.Unmarshal([]byte(rpctesting.BlindedDenebBlockContents), &signedblock)
			require.NoError(t, err)
			require.DeepEqual(t, converted, signedblock.SignedBlindedBlock.Message)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(rpctesting.BlindedDenebBlockContents)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlindedBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("invalid block", func(t *testing.T) {
		server := &Server{
			SyncChecker: &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(rpctesting.BellatrixBlock)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlindedBlockV2(writer, request)
		assert.Equal(t, http.StatusBadRequest, writer.Code)
		assert.Equal(t, true, strings.Contains(writer.Body.String(), "Body does not represent a valid block type"))
	})
	t.Run("syncing", func(t *testing.T) {
		chainService := &testing2.ChainService{}
		server := &Server{
			SyncChecker:           &mockSync.Sync{IsSyncing: true},
			HeadFetcher:           chainService,
			TimeFetcher:           chainService,
			OptimisticModeFetcher: chainService,
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte("foo")))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlindedBlockV2(writer, request)
		assert.Equal(t, http.StatusServiceUnavailable, writer.Code)
		assert.Equal(t, true, strings.Contains(writer.Body.String(), "Beacon node is currently syncing and not serving request on that endpoint"))
	})
}

func TestPublishBlindedBlockV2SSZ(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Run("Bellatrix", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			_, ok := req.Block.(*eth.GenericSignedBeaconBlock_BlindedBellatrix)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		var bellablock shared.SignedBlindedBeaconBlockBellatrix
		err := json.Unmarshal([]byte(rpctesting.BlindedBellatrixBlock), &bellablock)
		require.NoError(t, err)
		genericBlock, err := bellablock.ToGeneric()
		require.NoError(t, err)
		sszvalue, err := genericBlock.GetBlindedBellatrix().MarshalSSZ()
		require.NoError(t, err)
		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader(sszvalue))
		request.Header.Set("Accept", "application/octet-stream")
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlindedBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("Capella", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			_, ok := req.Block.(*eth.GenericSignedBeaconBlock_BlindedCapella)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		var cblock shared.SignedBlindedBeaconBlockCapella
		err := json.Unmarshal([]byte(rpctesting.BlindedCapellaBlock), &cblock)
		require.NoError(t, err)
		genericBlock, err := cblock.ToGeneric()
		require.NoError(t, err)
		sszvalue, err := genericBlock.GetBlindedCapella().MarshalSSZ()
		require.NoError(t, err)
		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader(sszvalue))
		request.Header.Set("Accept", "application/octet-stream")
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlindedBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("Deneb", func(t *testing.T) {
		v1alpha1Server := mock2.NewMockBeaconNodeValidatorServer(ctrl)
		v1alpha1Server.EXPECT().ProposeBeaconBlock(gomock.Any(), mock.MatchedBy(func(req *eth.GenericSignedBeaconBlock) bool {
			_, ok := req.Block.(*eth.GenericSignedBeaconBlock_BlindedDeneb)
			return ok
		}))
		server := &Server{
			V1Alpha1ValidatorServer: v1alpha1Server,
			SyncChecker:             &mockSync.Sync{IsSyncing: false},
		}

		var cblock shared.SignedBlindedBeaconBlockContentsDeneb
		err := json.Unmarshal([]byte(rpctesting.BlindedDenebBlockContents), &cblock)
		require.NoError(t, err)
		genericBlock, err := cblock.ToGeneric()
		require.NoError(t, err)
		v1block, err := migration.V1Alpha1SignedBlindedBlockAndBlobsDenebToV2Blinded(genericBlock.GetBlindedDeneb())
		require.NoError(t, err)
		sszvalue, err := v1block.MarshalSSZ()
		require.NoError(t, err)
		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader(sszvalue))
		request.Header.Set("Accept", "application/octet-stream")
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlindedBlockV2(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
	})
	t.Run("invalid block", func(t *testing.T) {
		server := &Server{
			SyncChecker: &mockSync.Sync{IsSyncing: false},
		}

		request := httptest.NewRequest(http.MethodPost, "http://foo.example", bytes.NewReader([]byte(rpctesting.BellatrixBlock)))
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}
		server.PublishBlindedBlockV2(writer, request)
		assert.Equal(t, http.StatusBadRequest, writer.Code)
		assert.Equal(t, true, strings.Contains(writer.Body.String(), "Body does not represent a valid block type"))
	})
}

func TestValidateConsensus(t *testing.T) {
	ctx := context.Background()

	parentState, privs := util.DeterministicGenesisState(t, params.MinimalSpecConfig().MinGenesisActiveValidatorCount)
	parentBlock, err := util.GenerateFullBlock(parentState, privs, util.DefaultBlockGenConfig(), parentState.Slot())
	require.NoError(t, err)
	parentSbb, err := blocks.NewSignedBeaconBlock(parentBlock)
	require.NoError(t, err)
	st, err := transition.ExecuteStateTransition(ctx, parentState, parentSbb)
	require.NoError(t, err)
	block, err := util.GenerateFullBlock(st, privs, util.DefaultBlockGenConfig(), st.Slot())
	require.NoError(t, err)
	sbb, err := blocks.NewSignedBeaconBlock(block)
	require.NoError(t, err)
	parentRoot, err := parentSbb.Block().HashTreeRoot()
	require.NoError(t, err)
	server := &Server{
		Blocker: &testutil.MockBlocker{RootBlockMap: map[[32]byte]interfaces.ReadOnlySignedBeaconBlock{parentRoot: parentSbb}},
		Stater:  &testutil.MockStater{StatesByRoot: map[[32]byte]state.BeaconState{bytesutil.ToBytes32(parentBlock.Block.StateRoot): parentState}},
	}

	require.NoError(t, server.validateConsensus(ctx, sbb))
}

func TestValidateEquivocation(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		st, err := util.NewBeaconState()
		require.NoError(t, err)
		require.NoError(t, st.SetSlot(10))
		fc := doublylinkedtree.New()
		require.NoError(t, fc.InsertNode(context.Background(), st, bytesutil.ToBytes32([]byte("root"))))
		server := &Server{
			ForkchoiceFetcher: &testing2.ChainService{ForkChoiceStore: fc},
		}
		blk, err := blocks.NewSignedBeaconBlock(util.NewBeaconBlock())
		require.NoError(t, err)
		blk.SetSlot(st.Slot() + 1)

		require.NoError(t, server.validateEquivocation(blk.Block()))
	})
	t.Run("block already exists", func(t *testing.T) {
		st, err := util.NewBeaconState()
		require.NoError(t, err)
		require.NoError(t, st.SetSlot(10))
		fc := doublylinkedtree.New()
		require.NoError(t, fc.InsertNode(context.Background(), st, bytesutil.ToBytes32([]byte("root"))))
		server := &Server{
			ForkchoiceFetcher: &testing2.ChainService{ForkChoiceStore: fc},
		}
		blk, err := blocks.NewSignedBeaconBlock(util.NewBeaconBlock())
		require.NoError(t, err)
		blk.SetSlot(st.Slot())

		assert.ErrorContains(t, "already exists", server.validateEquivocation(blk.Block()))
	})
}

func TestServer_GetBlockRoot(t *testing.T) {
	beaconDB := dbTest.SetupDB(t)
	ctx := context.Background()

	url := "http://example.com/eth/v1/beacon/blocks/{block_id}}/root"
	genBlk, blkContainers := fillDBTestBlocks(ctx, t, beaconDB)
	headBlock := blkContainers[len(blkContainers)-1]
	t.Run("get root", func(t *testing.T) {
		wsb, err := blocks.NewSignedBeaconBlock(headBlock.Block.(*eth.BeaconBlockContainer_Phase0Block).Phase0Block)
		require.NoError(t, err)

		mockChainFetcher := &testing2.ChainService{
			DB:                  beaconDB,
			Block:               wsb,
			Root:                headBlock.BlockRoot,
			FinalizedCheckPoint: &eth.Checkpoint{Root: blkContainers[64].BlockRoot},
			FinalizedRoots:      map[[32]byte]bool{},
		}

		bs := &Server{
			BeaconDB:              beaconDB,
			ChainInfoFetcher:      mockChainFetcher,
			HeadFetcher:           mockChainFetcher,
			OptimisticModeFetcher: mockChainFetcher,
			FinalizationFetcher:   mockChainFetcher,
		}

		root, err := genBlk.Block.HashTreeRoot()
		require.NoError(t, err)

		tests := []struct {
			name     string
			blockID  map[string]string
			want     string
			wantErr  string
			wantCode int
		}{
			{
				name:     "bad formatting",
				blockID:  map[string]string{"block_id": "3bad0"},
				wantErr:  "Could not parse block ID",
				wantCode: http.StatusBadRequest,
			},
			{
				name:     "canonical slot",
				blockID:  map[string]string{"block_id": "30"},
				want:     hexutil.Encode(blkContainers[30].BlockRoot),
				wantErr:  "",
				wantCode: http.StatusOK,
			},
			{
				name:     "head",
				blockID:  map[string]string{"block_id": "head"},
				want:     hexutil.Encode(headBlock.BlockRoot),
				wantErr:  "",
				wantCode: http.StatusOK,
			},
			{
				name:     "finalized",
				blockID:  map[string]string{"block_id": "finalized"},
				want:     hexutil.Encode(blkContainers[64].BlockRoot),
				wantErr:  "",
				wantCode: http.StatusOK,
			},
			{
				name:     "genesis",
				blockID:  map[string]string{"block_id": "genesis"},
				want:     hexutil.Encode(root[:]),
				wantErr:  "",
				wantCode: http.StatusOK,
			},
			{
				name:     "genesis root",
				blockID:  map[string]string{"block_id": hexutil.Encode(root[:])},
				want:     hexutil.Encode(root[:]),
				wantErr:  "",
				wantCode: http.StatusOK,
			},
			{
				name:     "root",
				blockID:  map[string]string{"block_id": hexutil.Encode(blkContainers[20].BlockRoot)},
				want:     hexutil.Encode(blkContainers[20].BlockRoot),
				wantErr:  "",
				wantCode: http.StatusOK,
			},
			{
				name:     "non-existent root",
				blockID:  map[string]string{"block_id": hexutil.Encode(bytesutil.PadTo([]byte("hi there"), 32))},
				wantErr:  "Could not find block",
				wantCode: http.StatusNotFound,
			},
			{
				name:     "slot",
				blockID:  map[string]string{"block_id": "40"},
				want:     hexutil.Encode(blkContainers[40].BlockRoot),
				wantErr:  "",
				wantCode: http.StatusOK,
			},
			{
				name:     "no block",
				blockID:  map[string]string{"block_id": "105"},
				wantErr:  "Could not find any blocks with given slot",
				wantCode: http.StatusNotFound,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				request := httptest.NewRequest(http.MethodGet, url, nil)
				request = mux.SetURLVars(request, tt.blockID)
				writer := httptest.NewRecorder()

				writer.Body = &bytes.Buffer{}

				bs.GetBlockRoot(writer, request)
				assert.Equal(t, tt.wantCode, writer.Code)
				resp := &BlockRootResponse{}
				require.NoError(t, json.Unmarshal(writer.Body.Bytes(), resp))
				if tt.wantErr != "" {
					require.ErrorContains(t, tt.wantErr, errors.New(writer.Body.String()))
					return
				}
				require.NotNil(t, resp)
				require.DeepEqual(t, resp.Data.Root, tt.want)
			})
		}
	})
	t.Run("execution optimistic", func(t *testing.T) {
		wsb, err := blocks.NewSignedBeaconBlock(headBlock.Block.(*eth.BeaconBlockContainer_Phase0Block).Phase0Block)
		require.NoError(t, err)

		mockChainFetcher := &testing2.ChainService{
			DB:                  beaconDB,
			Block:               wsb,
			Root:                headBlock.BlockRoot,
			FinalizedCheckPoint: &eth.Checkpoint{Root: blkContainers[64].BlockRoot},
			Optimistic:          true,
			FinalizedRoots:      map[[32]byte]bool{},
			OptimisticRoots: map[[32]byte]bool{
				bytesutil.ToBytes32(headBlock.BlockRoot): true,
			},
		}

		bs := &Server{
			BeaconDB:              beaconDB,
			ChainInfoFetcher:      mockChainFetcher,
			HeadFetcher:           mockChainFetcher,
			OptimisticModeFetcher: mockChainFetcher,
			FinalizationFetcher:   mockChainFetcher,
		}

		request := httptest.NewRequest(http.MethodGet, url, nil)
		request = mux.SetURLVars(request, map[string]string{"block_id": "head"})
		writer := httptest.NewRecorder()
		writer.Body = &bytes.Buffer{}

		bs.GetBlockRoot(writer, request)
		assert.Equal(t, http.StatusOK, writer.Code)
		resp := &BlockRootResponse{}
		require.NoError(t, json.Unmarshal(writer.Body.Bytes(), resp))
		require.DeepEqual(t, resp.ExecutionOptimistic, true)
	})
	t.Run("finalized", func(t *testing.T) {
		wsb, err := blocks.NewSignedBeaconBlock(headBlock.Block.(*eth.BeaconBlockContainer_Phase0Block).Phase0Block)
		require.NoError(t, err)

		mockChainFetcher := &testing2.ChainService{
			DB:                  beaconDB,
			Block:               wsb,
			Root:                headBlock.BlockRoot,
			FinalizedCheckPoint: &eth.Checkpoint{Root: blkContainers[64].BlockRoot},
			Optimistic:          true,
			FinalizedRoots: map[[32]byte]bool{
				bytesutil.ToBytes32(blkContainers[32].BlockRoot): true,
				bytesutil.ToBytes32(blkContainers[64].BlockRoot): false,
			},
		}

		bs := &Server{
			BeaconDB:              beaconDB,
			ChainInfoFetcher:      mockChainFetcher,
			HeadFetcher:           mockChainFetcher,
			OptimisticModeFetcher: mockChainFetcher,
			FinalizationFetcher:   mockChainFetcher,
		}
		t.Run("true", func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, url, nil)
			request = mux.SetURLVars(request, map[string]string{"block_id": "32"})
			writer := httptest.NewRecorder()
			writer.Body = &bytes.Buffer{}

			bs.GetBlockRoot(writer, request)
			assert.Equal(t, http.StatusOK, writer.Code)
			resp := &BlockRootResponse{}
			require.NoError(t, json.Unmarshal(writer.Body.Bytes(), resp))
			require.DeepEqual(t, resp.Finalized, true)
		})
		t.Run("false", func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, url, nil)
			request = mux.SetURLVars(request, map[string]string{"block_id": "64"})
			writer := httptest.NewRecorder()
			writer.Body = &bytes.Buffer{}

			bs.GetBlockRoot(writer, request)
			assert.Equal(t, http.StatusOK, writer.Code)
			resp := &BlockRootResponse{}
			require.NoError(t, json.Unmarshal(writer.Body.Bytes(), resp))
			require.DeepEqual(t, resp.Finalized, false)
		})
	})
}
