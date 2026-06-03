package schedule

import "testing"

func TestNames(t *testing.T) {
	if roomName("event", "date") != "event_date" {
		t.Fatal("unexpected room name")
	}
	if queueName("event", "date") != "queue:event_date" {
		t.Fatal("unexpected queue name")
	}
}

func TestAreaStat(t *testing.T) {
	user := "u1"
	stat := areaStat("a1", []Seat{{ID: "s1", ReservedUserID: &user}, {ID: "s2"}})
	if stat.AreaID != "a1" || stat.TotalSeatsNum != 2 || stat.ReservedSeatsNum != 1 {
		t.Fatalf("unexpected area stat: %#v", stat)
	}
}
