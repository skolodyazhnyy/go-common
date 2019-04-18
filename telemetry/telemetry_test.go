package telemetry

import (
	"context"
	"net"
	"testing"
	"time"
)

type noplog struct {
}

func (noplog) Warning(string, map[string]interface{}) {
}

type statsdServer struct {
	Addr string
	conn net.PacketConn
}

func (s *statsdServer) Stop() {
	s.conn.Close()
}

func (s *statsdServer) Read() (string, error) {
	buf := make([]byte, 1024)
	n, _, err := s.conn.ReadFrom(buf)
	if err != nil {
		return "", err
	}

	return string(buf[:n]), nil
}

func newTestStatsdServer(t *testing.T) *statsdServer {
	addr := "localhost:3114"
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		t.Fatal(err)
	}

	return &statsdServer{conn: conn, Addr: addr}
}

func TestTelemetry(t *testing.T) {
	srv := newTestStatsdServer(t)
	defer srv.Stop()

	tele, err := New(noplog{}, srv.Addr, "test.", nil)
	if err != nil {
		t.Fatal("Unable to connect to the statsdServer:", err)
	}

	defer tele.Stop()

	t.Run("prefix", func(t *testing.T) {
		tele.Gauge(context.Background(), "namespace", 10, nil)

		got, err := srv.Read()
		if err != nil {
			t.Fatal("Reading has failed:", err)
		}

		if got == "namespace:10.000000|g" {
			t.Errorf("Metric reported without prefix")
		}

		if want := "test.namespace:10.000000|g"; got != want {
			t.Errorf("UDP package does not match expected value, want: %#v, got %#v", want, got)
		}
	})

	t.Run("global tags", func(t *testing.T) {
		tele, err := New(noplog{}, srv.Addr, "test.", []string{"global1", "global2"})
		if err != nil {
			t.Fatal("Unable to connect to the statsd:", err)
		}

		defer tele.Stop()

		tele.Gauge(context.Background(), "tags", 10, []string{"tag1", "tag2"})

		got, err := srv.Read()
		if err != nil {
			t.Fatal("Reading has failed:", err)
		}

		if got == "test.tags:10.000000|g|#global1,global2" {
			t.Errorf("Metric reported only with global tags")
		}

		if got == "test.tags:10.000000|g|#tag1,tag2" {
			t.Errorf("Metric reported only with local tags")
		}

		if want := "test.tags:10.000000|g|#global1,global2,tag1,tag2"; got != want {
			t.Errorf("UDP package does not match expected value, want: %#v, got %#v", want, got)
		}
	})

	t.Run("context tags", func(t *testing.T) {
		ctx := AppendContext(context.Background(), []string{"ctx1"})
		ctx = AppendContext(ctx, []string{"ctx2"})

		tele.Gauge(ctx, "tags", 10, []string{"tag1", "tag2"})

		got, err := srv.Read()
		if err != nil {
			t.Fatal("Reading has failed:", err)
		}

		if got == "test.tags:10.000000|g|#ctx1,ctx2" {
			t.Errorf("Metric reported only with context tags")
		}

		if got == "test.tags:10.000000|g|#tag1,tag2" {
			t.Errorf("Metric reported only with local tags")
		}

		if want := "test.tags:10.000000|g|#ctx1,ctx2,tag1,tag2"; got != want {
			t.Errorf("UDP package does not match expected value, want: %#v, got %#v", want, got)
		}
	})

	t.Run("gauge", func(t *testing.T) {
		tele.Gauge(context.Background(), "gauge", 10, nil)

		got, err := srv.Read()
		if err != nil {
			t.Fatal("Reading has failed:", err)
		}

		if want := "test.gauge:10.000000|g"; got != want {
			t.Errorf("UDP package does not match expected value, want: %#v, got %#v", want, got)
		}
	})

	t.Run("timing", func(t *testing.T) {
		tele.Timing(context.Background(), "timing", time.Second, nil)

		got, err := srv.Read()
		if err != nil {
			t.Fatal("Reading has failed:", err)
		}

		if want := "test.timing:1000.000000|ms"; got != want {
			t.Errorf("UDP package does not match expected value, want: %#v, got %#v", want, got)
		}
	})

	t.Run("incr", func(t *testing.T) {
		tele.Incr(context.Background(), "incr", nil)

		got, err := srv.Read()
		if err != nil {
			t.Fatal("Reading has failed:", err)
		}

		if want := "test.incr:1|c"; got != want {
			t.Errorf("UDP package does not match expected value, want: %#v, got %#v", want, got)
		}
	})

	t.Run("incrby", func(t *testing.T) {
		tele.IncrBy(context.Background(), "incr_by", 31, nil)

		got, err := srv.Read()
		if err != nil {
			t.Fatal("Reading has failed:", err)
		}

		if want := "test.incr_by:31|c"; got != want {
			t.Errorf("UDP package does not match expected value, want: %#v, got %#v", want, got)
		}
	})

	t.Run("tags", func(t *testing.T) {
		tele.Incr(context.Background(), "tags", []string{"tag1", "tag2"})

		got, err := srv.Read()
		if err != nil {
			t.Fatal("Reading has failed:", err)
		}

		if want := "test.tags:1|c|#tag1,tag2"; got != want {
			t.Errorf("UDP package does not match expected value, want: %#v, got %#v", want, got)
		}
	})
}
