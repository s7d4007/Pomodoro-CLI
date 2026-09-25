package main

import "testing"

func testConfig() timerConfig {
	return timerConfig{
		focusSeconds:      2,
		shortBreakSeconds: 3,
		longBreakSeconds:  4,
		cycles:            2,
	}
}

func TestAdvancePhaseCompletesPomodoroCycle(t *testing.T) {
	model := initialModel(testConfig())

	model.advancePhase()
	if model.phase != shortBreak || model.cycle != 1 || model.remaining != 3 {
		t.Fatalf("after first focus: phase=%v cycle=%d remaining=%d", model.phase, model.cycle, model.remaining)
	}

	model.advancePhase()
	if model.phase != focus || model.cycle != 2 || model.remaining != 2 {
		t.Fatalf("after first break: phase=%v cycle=%d remaining=%d", model.phase, model.cycle, model.remaining)
	}

	model.advancePhase()
	if model.phase != longBreak || model.cycle != 2 || model.remaining != 4 {
		t.Fatalf("after final focus: phase=%v cycle=%d remaining=%d", model.phase, model.cycle, model.remaining)
	}

	model.advancePhase()
	if model.phase != focus || model.cycle != 1 || model.remaining != 2 {
		t.Fatalf("after long break: phase=%v cycle=%d remaining=%d", model.phase, model.cycle, model.remaining)
	}
}

func TestResetRestoresFirstFocusSession(t *testing.T) {
	model := initialModel(testConfig())
	model.advancePhase()
	model.running = false

	model.reset()

	if model.phase != focus || model.cycle != 1 || model.remaining != 2 || !model.running {
		t.Fatalf("reset produced phase=%v cycle=%d remaining=%d running=%v", model.phase, model.cycle, model.remaining, model.running)
	}
}

func TestTickAutomaticallyAdvancesPhase(t *testing.T) {
	config := testConfig()
	config.focusSeconds = 1
	timerModel := initialModel(config)

	updated, command := timerModel.Update(tickMsg{})
	updatedModel := updated.(model)

	if updatedModel.phase != shortBreak || updatedModel.remaining != config.shortBreakSeconds || !updatedModel.running {
		t.Fatalf("tick produced phase=%v remaining=%d running=%v", updatedModel.phase, updatedModel.remaining, updatedModel.running)
	}
	if command == nil {
		t.Fatal("expected timer to schedule the next tick")
	}
}
