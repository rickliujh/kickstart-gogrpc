package impl

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	pb "github.com/rickliujh/kickstart-gogrpc/pkg/api/grpc/pb/v1"
	"github.com/rickliujh/kickstart-gogrpc/pkg/service"
	. "github.com/rickliujh/kickstart-gogrpc/pkg/sql"
	. "github.com/rickliujh/kickstart-gogrpc/pkg/sql/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	anypb "google.golang.org/protobuf/types/known/anypb"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

var maxTestRunDuration = 180 * time.Second // 3 minutes

func TestScalar(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), maxTestRunDuration)
	defer cancel()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})).With(slog.String("env", "test"))

	ctrl := gomock.NewController(t)
	m := NewMockQuerier(ctrl)

	s := Server{
		counter:     service.Counter{},
		name:        "test-server",
		version:     "v0.0.1",
		environment: "test",
		db:          m,
		logger:      logger,
	}

	t.Run("scalar sans args", func(t *testing.T) {
		_, err := NewServer("", "", "", nil, logger)
		assert.Error(t, err)
		_, err = NewServer("test", "", "", nil, logger)
		assert.Error(t, err)
		_, err = NewServer("test", "test", "", nil, logger)
		assert.Error(t, err)
		_, err = NewServer("test", "", "test", nil, logger)
		assert.Error(t, err)
		_, err = NewServer("", "", "test", nil, logger)
		assert.Error(t, err)
	})

	t.Run("scalar sans args", func(t *testing.T) {
		assert.Panics(t, func() { s.Scalar(ctx, nil) }, "should panics if argument is not given")
	})

	t.Run("scalar with args", func(t *testing.T) {
		dict := map[string]any{"John": "Doe", "foo": "bar"}
		bs, _ := json.Marshal(dict)
		data := &anypb.Any{
			TypeUrl: "type.googleapis.com/json",
			Value:   bs,
		}
		req := &connect.Request[pb.ScalarRequest]{
			Msg: &pb.ScalarRequest{
				Content: &pb.Content{
					Id:   uuid.New().String(),
					Data: data,
				},
				Sent: tspb.Now(),
			},
		}

		m.EXPECT().
			ListAuthors(gomock.Any()).
			Return([]Author{{ID: 0, Name: "test name", Bio: pgtype.Text{}}}, nil).
			Times(1)

		// Scalar example
		reswrap, err := s.Scalar(ctx, req)
		res := reswrap.Msg
		if err != nil {
			t.Fatalf("error on scalar: %v", err)
		}

		assert.NotEmpty(t, res.GetRequestId())
		assert.Greater(t, res.GetMessageCount(), int64(0))
		assert.Equal(t, res.GetMessagesProcessed(), res.GetMessageCount())
		assert.Equal(t, success, res.GetProcessingDetails())
	})

	t.Run("stream sans args", func(t *testing.T) {
		assert.Panics(t, func() { s.Stream(ctx, nil) }, "should panics if argument is not given")
	})
}
