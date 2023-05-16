package notifier_test

import (
	"context"
	"sync"
	"testing"
	"time"

	ario "github.com/kahosan/aria2-rpc"
	"github.com/kahosan/aria2-rpc/internal/testutils"
	"github.com/kahosan/aria2-rpc/notifier"
)

func TestNotifyListener(t *testing.T) {

	t.Run("if notify is false, the listener will not be created", func(t *testing.T) {
		_, err := ario.NewClient(testutils.Arai2Uri("https://"), "", false)
		if err != nil {
			t.Fatal("should error")
		}
	})

	client, err := ario.NewClient(testutils.Arai2Uri("https://"), "", true)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	t.Run("if notify is true, the listener will be created", func(t *testing.T) {
		notify, err := client.NotifyListener(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		defer notify.Close()

		gid, err := client.AddURI([]string{"https://releases.ubuntu.com/22.04.2/ubuntu-22.04.2-live-server-amd64.iso"}, nil)
		if err != nil {
			t.Fatal(err)
		}

		var wg sync.WaitGroup
		wg.Add(1)
		// using coroutines to prevent blocking
		go func() {
			for v := range notify.Start() {
				if v == gid {
					t.Log("task start: ", v)
				}
			}
		}()

		go func() {
			for v := range notify.Pause() {
				if v == gid {
					t.Log("task pause: ", v)
				}
			}
		}()

		go func() {
			for v := range notify.Stop() {
				if v == gid {
					t.Log("task stop: ", v)
					wg.Done()
				}
			}
		}()

		time.Sleep(time.Second)
		err = client.Pause(gid)
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Second)
		err = client.Unpause(gid)
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Second)
		err = client.Remove(gid)
		if err != nil {
			t.Fatal(err)
		}

		wg.Wait()
	})

	t.Run("multiple Listen once tests", func(t *testing.T) {
		notify, err := client.NotifyListener(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		defer notify.Close()

		gid, err := client.AddURI([]string{"https://releases.ubuntu.com/22.04.2/ubuntu-22.04.2-live-server-amd64.iso"}, nil)
		if err != nil {
			t.Fatal(err)
		}

		var wg sync.WaitGroup
		wg.Add(1)
		go notify.ListenOnce(notifier.NotifyEvents.Start, func(g string, stop func()) {
			if g == gid {
				t.Log("start: ", gid)
			}
		})

		go notify.ListenOnce(notifier.NotifyEvents.Pause, func(g string, stop func()) {
			if g == gid {
				t.Log("pause: ", gid)
			}
		})

		go notify.ListenOnce(notifier.NotifyEvents.Stop, func(g string, stop func()) {
			if g == gid {
				t.Log("stop: ", gid)
			}
			wg.Done()
		})

		time.Sleep(time.Second)
		err = client.Pause(gid)
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Second)
		err = client.Unpause(gid)
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Second)
		err = client.Remove(gid)
		if err != nil {
			t.Fatal(err)
		}

		wg.Wait()
	})

	t.Run("multiple Listen tests", func(t *testing.T) {
		notify, err := client.NotifyListener(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		defer notify.Close()

		gid, err := client.AddURI([]string{"https://releases.ubuntu.com/22.04.2/ubuntu-22.04.2-live-server-amd64.iso"}, nil)
		if err != nil {
			t.Fatal(err)
		}

		var wg = sync.WaitGroup{}
		wg.Add(1)

		task := notifier.Tasks{
			notifier.NotifyEvents.Start: func(g string) {
				if g == gid {
					t.Log("下载开始: ", g)
				}
			},
			notifier.NotifyEvents.Complete: func(g string) {
				if g == gid {
					t.Log("下载完成: ", g)
					wg.Done()
				}
			},
			notifier.NotifyEvents.Error: func(g string) {
				if g == gid {
					t.Log("下载失败: ", g)
					wg.Done()
				}
			},
			notifier.NotifyEvents.Stop: func(g string) {
				if g == gid {
					t.Log("下载移除: ", g)
					wg.Done()
				}
			},
		}

		notify.ListenMultiple(task)

		time.Sleep(time.Second)
		err = client.Remove(gid)
		if err != nil {
			t.Fatal(err)
		}

		wg.Wait()
	})
}
