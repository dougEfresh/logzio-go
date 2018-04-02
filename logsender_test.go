// Copyright © 2017 Douglas Chimento <dchimento@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logzio

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestLogzioSender_Retries(t *testing.T) {
	var sent= make([]byte, 1024)
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		r.Body.Read(sent)
	}))
	defer ts.Close()
	l, err := New(
		"fake-token",
		SetDebug(os.Stderr),
		SetUrl("http://localhost:12345"),
		SetDrainDuration(time.Minute*10),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Stop()
	l.Send([]byte("blah"))
	l.Drain()
	item, err := l.queue.Dequeue()
	// expected msg to be in queue after max retries
	if item == nil || item.ID != 2{
		t.Fatalf("Unexpect item in the queue - %s",string(item.Value))
	}
	item, err = l.queue.Dequeue()
	// expected queue to be empty - only one requeue executed
	if err == nil{
		t.Fatalf("Unexpect item in the queue - %s",string(item.Value))
	}
}

func TestLogzioSender_Send(t *testing.T) {
	var sent = make([]byte, 1024)
	var sentToken string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sentToken = r.URL.Query().Get("token")
		w.WriteHeader(http.StatusOK)
		r.Body.Read(sent)
	}))
	defer ts.Close()

	l, err := New("fake-token", SetUrl(ts.URL))
	defer l.Stop()
	if err != nil {
		t.Fatal(err)
	}
	l.Send([]byte("blah"))
	l.Drain()
	time.Sleep(200 * time.Millisecond)
	sentMsg := string(sent[0:5])
	if sentMsg != "blah\n" {
		t.Fatalf("%s != %s ", sent, sentMsg)
	}
	if sentToken != "fake-token" {
		t.Fatalf("token not sent %s", sentToken)
	}
}

func TestLogzioSender_DelayStart(t *testing.T) {
	var sent = make([]byte, 1024)
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		r.Body.Read(sent)
	}))
	defer ts.Close()
	l, err := New(
		"fake-token",
		SetDebug(os.Stderr),
		SetUrl("http://localhost:12345"),
	)
	if err != nil {
		t.Fatal(err)
	}
	l.Send([]byte("blah"))
	time.Sleep(200 * time.Millisecond)
	l.Drain()
	ts.Start()
	SetUrl(ts.URL)(l)
	l.Drain()
	time.Sleep(500 * time.Millisecond)
	sentMsg := string(sent[0:5])
	if len(sentMsg) != 5 {
		t.Fatalf("Wrong len of msg %d", len(sentMsg))
	}
	if sentMsg != "blah\n" {
		t.Fatalf("%s != %s ", sent, sentMsg)
	}
}

func TestLogzioSender_TmpDir(t *testing.T) {
	var sent = make([]byte, 1024)
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		r.Body.Read(sent)
	}))
	ts.Start()
	defer ts.Close()
	tmp := fmt.Sprintf("%s/%d", os.TempDir(), time.Now().Nanosecond())
	l, err := New(
		"fake-token",
		SetDebug(os.Stderr),
		SetTempDirectory(tmp),
		SetDrainDuration(time.Minute),
		SetUrl(ts.URL),
	)
	if err != nil {
		t.Fatal(err)
	}
	l.Send([]byte("blah"))
	time.Sleep(200 * time.Millisecond)
	l.Drain()
	sentMsg := string(sent[0:5])
	if len(sentMsg) != 5 {
		t.Fatalf("Wrong len of msg %d", len(sentMsg))
	}
	if sentMsg != "blah\n" {
		t.Fatalf("%s != %s ", string(sent), string(sentMsg))
	}
}

func TestLogzioSender_Write(t *testing.T) {
	var sent = make([]byte, 1024)
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		r.Body.Read(sent)
	}))
	ts.Start()
	defer ts.Close()
	tmp := fmt.Sprintf("%s/%d", os.TempDir(), time.Now().Nanosecond())
	l, err := New(
		"fake-token",
		SetDebug(os.Stderr),
		SetTempDirectory(tmp),
		SetDrainDuration(time.Minute),
		SetUrl(ts.URL),
	)
	if err != nil {
		t.Fatal(err)
	}
	l.Write([]byte("blah"))
	time.Sleep(200 * time.Millisecond)
	l.Sync()
	sentMsg := string(sent[0:5])
	if len(sentMsg) != 5 {
		t.Fatalf("Wrong len of msg %d", len(sentMsg))
	}
	if sentMsg != "blah\n" {
		t.Fatalf("%s != %s ", string(sent), string(sentMsg))
	}
}

func TestLogzioSender_Unauth(t *testing.T) {
	var sent = make([]byte, 1024)
	cnt := 0
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cnt++
		if cnt == 2 {
			w.WriteHeader(http.StatusAccepted)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
		r.Body.Read(sent)
	}))
	ts.Start()
	defer ts.Close()
	tmp := fmt.Sprintf("%s/%d", os.TempDir(), time.Now().Nanosecond())
	l, err := New(
		"fake-token",
		SetDebug(os.Stderr),
		SetTempDirectory(tmp),
		SetDrainDuration(time.Minute),
		SetUrl(ts.URL),
	)
	if err != nil {
		t.Fatal(err)
	}
	l.Write([]byte("blah"))
	time.Sleep(200 * time.Millisecond)
	l.Sync()
	time.Sleep(100 * time.Millisecond)
	l.Drain()
	time.Sleep(100 * time.Millisecond)
	sentMsg := string(sent[0:5])
	if len(sentMsg) != 5 {
		t.Fatalf("Wrong len of msg %d", len(sentMsg))
	}
	if sentMsg != "blah\n" {
		t.Fatalf("%s != %s ", string(sent), string(sentMsg))
	}
}

func TestLogzioSender_ThresholdLimit(t *testing.T) {
	l, err := New(
		"fake-token",
		SetDebug(os.Stderr),
		SetUrl("http://localhost:12345"),
		SetDrainDiskThreshold(0),
		SetDrainDuration(time.Minute),
	)
	if err != nil {
		t.Fatal(err)
	}

	l.Send([]byte("blah"))
	item, err := l.queue.Dequeue()
	if item != nil{
		t.Fatalf("Unexpect item in the queue - %s",string(item.Value))
	}
}

func TestLogzioSender_ThresholdLimitWithoutCheck(t *testing.T) {
	l, err := New(
		"fake-token",
		SetDebug(os.Stderr),
		SetUrl("http://localhost:12345"),
		SetDrainDiskThreshold(0),
		SetCheckDiskSpace(false),
		SetDrainDuration(time.Minute),
	)
	if err != nil {
		t.Fatal(err)
	}

	l.Send([]byte("blah"))
	item, err := l.queue.Dequeue()
	if item == nil{
		t.Fatalf("Unexpect item in the queue - %s",string(item.Value))
	}
}

func TestLogzioSender_RestoreQueue(t *testing.T) {
	var sent= make([]byte, 1024)
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		r.Body.Read(sent)
	}))
	defer ts.Close()
	l, err := New(
		"fake-token",
		SetDebug(os.Stderr),
		SetUrl("http://localhost:12345"),
		SetDrainDuration(time.Minute*10),
		SetTempDirectory("./data"),
	)
	if err != nil {
		t.Fatal(err)
	}
	l.Send([]byte("blah"))
	l.Stop()

	// open queue again - same dir
	l, err = New(
		"fake-token",
		SetDebug(os.Stderr),
		SetUrl("http://localhost:12345"),
		SetDrainDuration(time.Minute*10),
		SetTempDirectory("./data"),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Stop()

	item, err := l.queue.Dequeue()
	if string(item.Value) != "blah\n"{
		t.Fatalf("Unexpect item in the queue - %s", string(item.Value))
	}
	if item.ID != 2{
		t.Fatalf("Unexpect ID number - %s", string(item.ID))
	}
}

func BenchmarkLogzioSender(b *testing.B) {
	b.ReportAllocs()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	l, _ := New("fake-token", SetUrl(ts.URL), SetDrainDuration(time.Hour))
	defer ts.Close()
	defer l.Stop()
	msg := []byte("test")
	for i := 0; i < b.N; i++ {
		l.Send(msg)
	}
}
