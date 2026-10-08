package repository

import (
	"testing"

	"darbelis.eu/taxi/internal/dto"
)

func TestChatMessageRepositoryMemory_AddRemoveScenario(t *testing.T) {
	repo := NewChatMessageRepositoryMemory()

	ab1 := dto.NewChatMessage("ab1", "alice", "bob", "hi bob")
	ab2 := dto.NewChatMessage("ab2", "bob", "alice", "hi alice")
	ab3 := dto.NewChatMessage("ab3", "alice", "bob", "how are you")
	ac1 := dto.NewChatMessage("ac1", "alice", "carol", "hi carol")

	repo.Add(ab1)
	repo.Add(ab2)
	repo.Add(ab3)
	repo.Add(ac1)

	if removed := repo.Remove("ab2"); !removed {
		t.Fatalf("Remove(ab2) = false, want true")
	}

	if removed := repo.Remove("missing"); removed {
		t.Fatalf("Remove(missing) = true, want false")
	}

	if got := repo.Get("ab2"); got != nil {
		t.Errorf("Get(ab2) after removal = %v, want nil", got)
	}

	if got := repo.Get("ab1"); got != ab1 {
		t.Errorf("Get(ab1) = %v, want %v", got, ab1)
	}

	abConversation := repo.GetConversation("alice", "bob")
	wantAB := []*dto.ChatMessage{ab1, ab3}
	assertConversationEquals(t, abConversation, wantAB)

	// Pair key is order-independent.
	abConversationReversed := repo.GetConversation("bob", "alice")
	assertConversationEquals(t, abConversationReversed, wantAB)

	acConversation := repo.GetConversation("alice", "carol")
	assertConversationEquals(t, acConversation, []*dto.ChatMessage{ac1})

	// A partner should be listed once per distinct conversation partner, not once
	// per message exchanged with them.
	assertStringSliceEquals(t, repo.GetChatPartners("alice"), []string{"bob", "carol"})
	assertStringSliceEquals(t, repo.GetChatPartners("bob"), []string{"alice"})
	assertStringSliceEquals(t, repo.GetChatPartners("carol"), []string{"alice"})
	assertStringSliceEquals(t, repo.GetChatPartners("dave"), []string{})

	repo.Remove("ac1")
	assertStringSliceEquals(t, repo.GetChatPartners("carol"), []string{})

}

func assertConversationEquals(t *testing.T, got, want []*dto.ChatMessage) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("conversation length = %d, want %d (got %v)", len(got), len(want), got)
	}

	for i := range want {
		if got[i].Id != want[i].Id {
			t.Errorf("conversation[%d].Id = %q, want %q", i, got[i].Id, want[i].Id)
		}
	}
}

func assertStringSliceEquals(t *testing.T, got, want []string) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("slice length = %d, want %d (got %v)", len(got), len(want), got)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf("[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
