// Package main implements a personal life management MCP server.
// It provides tools for managing tasks, notes, journal entries,
// shopping lists, and reminders — all stored in a local JSON file.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ---- Data model ------------------------------------------------------------

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

type Task struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Done      bool     `json:"done"`
	Priority  Priority `json:"priority"`
	DueDate   string   `json:"due_date,omitempty"`
	CreatedAt string   `json:"created_at"`
}

type Note struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Tags      []string `json:"tags,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type JournalEntry struct {
	ID        string `json:"id"`
	Mood      string `json:"mood,omitempty"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type ShoppingItem struct {
	ID        string `json:"id"`
	Item      string `json:"item"`
	Quantity  string `json:"quantity,omitempty"`
	Checked   bool   `json:"checked"`
	CreatedAt string `json:"created_at"`
}

type Reminder struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	When        string `json:"when"`
	Dismissed   bool   `json:"dismissed"`
	CreatedAt   string `json:"created_at"`
}

type Store struct {
	Tasks         []Task         `json:"tasks"`
	Notes         []Note         `json:"notes"`
	JournalEntries []JournalEntry `json:"journal_entries"`
	ShoppingItems []ShoppingItem  `json:"shopping_items"`
	Reminders     []Reminder     `json:"reminders"`
}

// ---- Storage ---------------------------------------------------------------

type DB struct {
	path  string
	store Store
}

func newDB(path string) (*DB, error) {
	db := &DB{path: path}
	if err := db.load(); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) load() error {
	data, err := os.ReadFile(db.path)
	if os.IsNotExist(err) {
		db.store = Store{}
		return nil
	}
	if err != nil {
		return fmt.Errorf("reading store: %w", err)
	}
	return json.Unmarshal(data, &db.store)
}

func (db *DB) save() error {
	if err := os.MkdirAll(filepath.Dir(db.path), 0o700); err != nil {
		return fmt.Errorf("creating data dir: %w", err)
	}
	data, err := json.MarshalIndent(db.store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling store: %w", err)
	}
	return os.WriteFile(db.path, data, 0o600)
}

func genID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func now() string {
	return time.Now().Format(time.RFC3339)
}

// ---- Tool result helpers ---------------------------------------------------

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}
}

func errResult(msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
		IsError: true,
	}
}

// ---- Argument parsing helper -----------------------------------------------

func parseArgs(req *mcp.CallToolRequest, v any) error {
	if req.Params.Arguments == nil {
		return fmt.Errorf("missing arguments")
	}
	b, err := json.Marshal(req.Params.Arguments)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

// ---- Server ----------------------------------------------------------------

type Server struct {
	db *DB
}

// task_add
func (s *Server) taskAdd(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Title    string   `json:"title"`
		Priority Priority `json:"priority"`
		DueDate  string   `json:"due_date"`
	}
	if err := parseArgs(req, &args); err != nil {
		return errResult(err.Error()), nil
	}
	if args.Title == "" {
		return errResult("title is required"), nil
	}
	if args.Priority == "" {
		args.Priority = PriorityMedium
	}
	task := Task{
		ID:        genID(),
		Title:     args.Title,
		Priority:  args.Priority,
		DueDate:   args.DueDate,
		CreatedAt: now(),
	}
	s.db.store.Tasks = append(s.db.store.Tasks, task)
	if err := s.db.save(); err != nil {
		return errResult(fmt.Sprintf("failed to save: %v", err)), nil
	}
	return textResult(fmt.Sprintf("Added task [%s]: %s", task.ID, task.Title)), nil
}

// task_list
func (s *Server) taskList(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		ShowDone bool `json:"show_done"`
	}
	_ = parseArgs(req, &args)

	var lines []string
	for _, t := range s.db.store.Tasks {
		if t.Done && !args.ShowDone {
			continue
		}
		status := "[ ]"
		if t.Done {
			status = "[x]"
		}
		due := ""
		if t.DueDate != "" {
			due = fmt.Sprintf(" (due: %s)", t.DueDate)
		}
		lines = append(lines, fmt.Sprintf("%s [%s] %s%s  (id: %s)", status, strings.ToUpper(string(t.Priority)), t.Title, due, t.ID))
	}
	if len(lines) == 0 {
		return textResult("No tasks found."), nil
	}
	return textResult(strings.Join(lines, "\n")), nil
}

// task_complete
func (s *Server) taskComplete(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		ID string `json:"id"`
	}
	if err := parseArgs(req, &args); err != nil {
		return errResult(err.Error()), nil
	}
	for i, t := range s.db.store.Tasks {
		if t.ID == args.ID {
			s.db.store.Tasks[i].Done = true
			if err := s.db.save(); err != nil {
				return errResult(fmt.Sprintf("failed to save: %v", err)), nil
			}
			return textResult(fmt.Sprintf("Completed task: %s", t.Title)), nil
		}
	}
	return errResult(fmt.Sprintf("task not found: %s", args.ID)), nil
}

// task_delete
func (s *Server) taskDelete(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		ID string `json:"id"`
	}
	if err := parseArgs(req, &args); err != nil {
		return errResult(err.Error()), nil
	}
	for i, t := range s.db.store.Tasks {
		if t.ID == args.ID {
			s.db.store.Tasks = append(s.db.store.Tasks[:i], s.db.store.Tasks[i+1:]...)
			if err := s.db.save(); err != nil {
				return errResult(fmt.Sprintf("failed to save: %v", err)), nil
			}
			return textResult(fmt.Sprintf("Deleted task: %s", t.Title)), nil
		}
	}
	return errResult(fmt.Sprintf("task not found: %s", args.ID)), nil
}

// note_add
func (s *Server) noteAdd(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Title string   `json:"title"`
		Body  string   `json:"body"`
		Tags  []string `json:"tags"`
	}
	if err := parseArgs(req, &args); err != nil {
		return errResult(err.Error()), nil
	}
	if args.Title == "" {
		return errResult("title is required"), nil
	}
	n := Note{
		ID:        genID(),
		Title:     args.Title,
		Body:      args.Body,
		Tags:      args.Tags,
		CreatedAt: now(),
		UpdatedAt: now(),
	}
	s.db.store.Notes = append(s.db.store.Notes, n)
	if err := s.db.save(); err != nil {
		return errResult(fmt.Sprintf("failed to save: %v", err)), nil
	}
	return textResult(fmt.Sprintf("Added note [%s]: %s", n.ID, n.Title)), nil
}

// note_list
func (s *Server) noteList(_ context.Context, _ *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if len(s.db.store.Notes) == 0 {
		return textResult("No notes found."), nil
	}
	var lines []string
	for _, n := range s.db.store.Notes {
		tags := ""
		if len(n.Tags) > 0 {
			tags = fmt.Sprintf(" [%s]", strings.Join(n.Tags, ", "))
		}
		lines = append(lines, fmt.Sprintf("• %s%s  (id: %s)", n.Title, tags, n.ID))
	}
	return textResult(strings.Join(lines, "\n")), nil
}

// note_get
func (s *Server) noteGet(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		ID string `json:"id"`
	}
	if err := parseArgs(req, &args); err != nil {
		return errResult(err.Error()), nil
	}
	for _, n := range s.db.store.Notes {
		if n.ID == args.ID {
			return textResult(fmt.Sprintf("# %s\nTags: %s\nCreated: %s\n\n%s", n.Title, strings.Join(n.Tags, ", "), n.CreatedAt, n.Body)), nil
		}
	}
	return errResult(fmt.Sprintf("note not found: %s", args.ID)), nil
}

// note_search
func (s *Server) noteSearch(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Query string `json:"query"`
	}
	if err := parseArgs(req, &args); err != nil {
		return errResult(err.Error()), nil
	}
	query := strings.ToLower(args.Query)
	var lines []string
	for _, n := range s.db.store.Notes {
		if strings.Contains(strings.ToLower(n.Title), query) ||
			strings.Contains(strings.ToLower(n.Body), query) {
			lines = append(lines, fmt.Sprintf("• %s  (id: %s)", n.Title, n.ID))
		}
	}
	if len(lines) == 0 {
		return textResult("No matching notes."), nil
	}
	return textResult(strings.Join(lines, "\n")), nil
}

// note_delete
func (s *Server) noteDelete(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		ID string `json:"id"`
	}
	if err := parseArgs(req, &args); err != nil {
		return errResult(err.Error()), nil
	}
	for i, n := range s.db.store.Notes {
		if n.ID == args.ID {
			s.db.store.Notes = append(s.db.store.Notes[:i], s.db.store.Notes[i+1:]...)
			if err := s.db.save(); err != nil {
				return errResult(fmt.Sprintf("failed to save: %v", err)), nil
			}
			return textResult(fmt.Sprintf("Deleted note: %s", n.Title)), nil
		}
	}
	return errResult(fmt.Sprintf("note not found: %s", args.ID)), nil
}

// journal_add
func (s *Server) journalAdd(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Content string `json:"content"`
		Mood    string `json:"mood"`
	}
	if err := parseArgs(req, &args); err != nil {
		return errResult(err.Error()), nil
	}
	if args.Content == "" {
		return errResult("content is required"), nil
	}
	entry := JournalEntry{
		ID:        genID(),
		Mood:      args.Mood,
		Content:   args.Content,
		CreatedAt: now(),
	}
	s.db.store.JournalEntries = append(s.db.store.JournalEntries, entry)
	if err := s.db.save(); err != nil {
		return errResult(fmt.Sprintf("failed to save: %v", err)), nil
	}
	return textResult(fmt.Sprintf("Journal entry saved [%s] at %s", entry.ID, entry.CreatedAt)), nil
}

// journal_list
func (s *Server) journalList(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Limit int `json:"limit"`
	}
	_ = parseArgs(req, &args)
	if args.Limit <= 0 {
		args.Limit = 10
	}

	entries := make([]JournalEntry, len(s.db.store.JournalEntries))
	copy(entries, s.db.store.JournalEntries)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].CreatedAt > entries[j].CreatedAt
	})
	if len(entries) > args.Limit {
		entries = entries[:args.Limit]
	}
	if len(entries) == 0 {
		return textResult("No journal entries found."), nil
	}
	var parts []string
	for _, e := range entries {
		mood := ""
		if e.Mood != "" {
			mood = fmt.Sprintf(" (mood: %s)", e.Mood)
		}
		parts = append(parts, fmt.Sprintf("--- %s%s ---\n%s", e.CreatedAt, mood, e.Content))
	}
	return textResult(strings.Join(parts, "\n\n")), nil
}

// shopping_add
func (s *Server) shoppingAdd(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Item     string `json:"item"`
		Quantity string `json:"quantity"`
	}
	if err := parseArgs(req, &args); err != nil {
		return errResult(err.Error()), nil
	}
	if args.Item == "" {
		return errResult("item is required"), nil
	}
	si := ShoppingItem{
		ID:        genID(),
		Item:      args.Item,
		Quantity:  args.Quantity,
		CreatedAt: now(),
	}
	s.db.store.ShoppingItems = append(s.db.store.ShoppingItems, si)
	if err := s.db.save(); err != nil {
		return errResult(fmt.Sprintf("failed to save: %v", err)), nil
	}
	return textResult(fmt.Sprintf("Added to shopping list: %s", si.Item)), nil
}

// shopping_list
func (s *Server) shoppingList(_ context.Context, _ *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if len(s.db.store.ShoppingItems) == 0 {
		return textResult("Shopping list is empty."), nil
	}
	var lines []string
	for _, si := range s.db.store.ShoppingItems {
		check := "[ ]"
		if si.Checked {
			check = "[x]"
		}
		qty := ""
		if si.Quantity != "" {
			qty = fmt.Sprintf(" (%s)", si.Quantity)
		}
		lines = append(lines, fmt.Sprintf("%s %s%s  (id: %s)", check, si.Item, qty, si.ID))
	}
	return textResult(strings.Join(lines, "\n")), nil
}

// shopping_check
func (s *Server) shoppingCheck(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		ID string `json:"id"`
	}
	if err := parseArgs(req, &args); err != nil {
		return errResult(err.Error()), nil
	}
	for i, si := range s.db.store.ShoppingItems {
		if si.ID == args.ID {
			s.db.store.ShoppingItems[i].Checked = true
			if err := s.db.save(); err != nil {
				return errResult(fmt.Sprintf("failed to save: %v", err)), nil
			}
			return textResult(fmt.Sprintf("Checked off: %s", si.Item)), nil
		}
	}
	return errResult(fmt.Sprintf("shopping item not found: %s", args.ID)), nil
}

// shopping_clear
func (s *Server) shoppingClear(_ context.Context, _ *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.db.store.ShoppingItems = nil
	if err := s.db.save(); err != nil {
		return errResult(fmt.Sprintf("failed to save: %v", err)), nil
	}
	return textResult("Shopping list cleared."), nil
}

// reminder_add
func (s *Server) reminderAdd(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		Title string `json:"title"`
		When  string `json:"when"`
	}
	if err := parseArgs(req, &args); err != nil {
		return errResult(err.Error()), nil
	}
	if args.Title == "" {
		return errResult("title is required"), nil
	}
	r := Reminder{
		ID:        genID(),
		Title:     args.Title,
		When:      args.When,
		CreatedAt: now(),
	}
	s.db.store.Reminders = append(s.db.store.Reminders, r)
	if err := s.db.save(); err != nil {
		return errResult(fmt.Sprintf("failed to save: %v", err)), nil
	}
	return textResult(fmt.Sprintf("Reminder set [%s]: %s", r.ID, r.Title)), nil
}

// reminder_list
func (s *Server) reminderList(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		ShowDismissed bool `json:"show_dismissed"`
	}
	_ = parseArgs(req, &args)

	var lines []string
	for _, r := range s.db.store.Reminders {
		if r.Dismissed && !args.ShowDismissed {
			continue
		}
		when := ""
		if r.When != "" {
			when = fmt.Sprintf(" @ %s", r.When)
		}
		status := ""
		if r.Dismissed {
			status = " [dismissed]"
		}
		lines = append(lines, fmt.Sprintf("• %s%s%s  (id: %s)", r.Title, when, status, r.ID))
	}
	if len(lines) == 0 {
		return textResult("No active reminders."), nil
	}
	return textResult(strings.Join(lines, "\n")), nil
}

// reminder_dismiss
func (s *Server) reminderDismiss(_ context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args struct {
		ID string `json:"id"`
	}
	if err := parseArgs(req, &args); err != nil {
		return errResult(err.Error()), nil
	}
	for i, r := range s.db.store.Reminders {
		if r.ID == args.ID {
			s.db.store.Reminders[i].Dismissed = true
			if err := s.db.save(); err != nil {
				return errResult(fmt.Sprintf("failed to save: %v", err)), nil
			}
			return textResult(fmt.Sprintf("Dismissed reminder: %s", r.Title)), nil
		}
	}
	return errResult(fmt.Sprintf("reminder not found: %s", args.ID)), nil
}

// ---- Tool definitions ------------------------------------------------------

func strProp(desc string) *jsonschema.Schema {
	return &jsonschema.Schema{Type: "string", Description: desc}
}

func boolProp(desc string) *jsonschema.Schema {
	return &jsonschema.Schema{Type: "boolean", Description: desc}
}

func intProp(desc string) *jsonschema.Schema {
	return &jsonschema.Schema{Type: "integer", Description: desc}
}

func strArrProp(desc string) *jsonschema.Schema {
	return &jsonschema.Schema{Type: "array", Description: desc, Items: &jsonschema.Schema{Type: "string"}}
}

func buildTools(srv *Server) []*mcp.Tool {
	return []*mcp.Tool{
		{
			Name:        "task_add",
			Description: "Add a new task to your to-do list",
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"title":    strProp("Task description"),
					"priority": {Type: "string", Description: "Priority level", Enum: []any{"low", "medium", "high"}},
					"due_date": strProp("Optional due date (e.g. 2025-12-31)"),
				},
				Required: []string{"title"},
			},
		},
		{
			Name:        "task_list",
			Description: "List your tasks",
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"show_done": boolProp("Include completed tasks (default: false)"),
				},
			},
		},
		{
			Name:        "task_complete",
			Description: "Mark a task as done",
			InputSchema: &jsonschema.Schema{
				Type:       "object",
				Properties: map[string]*jsonschema.Schema{"id": strProp("Task ID")},
				Required:   []string{"id"},
			},
		},
		{
			Name:        "task_delete",
			Description: "Delete a task",
			InputSchema: &jsonschema.Schema{
				Type:       "object",
				Properties: map[string]*jsonschema.Schema{"id": strProp("Task ID")},
				Required:   []string{"id"},
			},
		},
		{
			Name:        "note_add",
			Description: "Save a new note",
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"title": strProp("Note title"),
					"body":  strProp("Note content"),
					"tags":  strArrProp("Optional tags"),
				},
				Required: []string{"title"},
			},
		},
		{
			Name:        "note_list",
			Description: "List all notes",
			InputSchema: &jsonschema.Schema{Type: "object"},
		},
		{
			Name:        "note_get",
			Description: "Get the full content of a note by ID",
			InputSchema: &jsonschema.Schema{
				Type:       "object",
				Properties: map[string]*jsonschema.Schema{"id": strProp("Note ID")},
				Required:   []string{"id"},
			},
		},
		{
			Name:        "note_search",
			Description: "Search notes by keyword",
			InputSchema: &jsonschema.Schema{
				Type:       "object",
				Properties: map[string]*jsonschema.Schema{"query": strProp("Search keyword")},
				Required:   []string{"query"},
			},
		},
		{
			Name:        "note_delete",
			Description: "Delete a note",
			InputSchema: &jsonschema.Schema{
				Type:       "object",
				Properties: map[string]*jsonschema.Schema{"id": strProp("Note ID")},
				Required:   []string{"id"},
			},
		},
		{
			Name:        "journal_add",
			Description: "Add a journal entry",
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"content": strProp("Journal entry text"),
					"mood":    strProp("Optional mood tag (e.g. happy, tired, anxious)"),
				},
				Required: []string{"content"},
			},
		},
		{
			Name:        "journal_list",
			Description: "List recent journal entries",
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"limit": intProp("Max entries to return (default: 10)"),
				},
			},
		},
		{
			Name:        "shopping_add",
			Description: "Add an item to the shopping list",
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"item":     strProp("Item name"),
					"quantity": strProp("Optional quantity (e.g. 2, 1kg)"),
				},
				Required: []string{"item"},
			},
		},
		{
			Name:        "shopping_list",
			Description: "View the shopping list",
			InputSchema: &jsonschema.Schema{Type: "object"},
		},
		{
			Name:        "shopping_check",
			Description: "Mark a shopping item as picked up",
			InputSchema: &jsonschema.Schema{
				Type:       "object",
				Properties: map[string]*jsonschema.Schema{"id": strProp("Shopping item ID")},
				Required:   []string{"id"},
			},
		},
		{
			Name:        "shopping_clear",
			Description: "Clear the entire shopping list",
			InputSchema: &jsonschema.Schema{Type: "object"},
		},
		{
			Name:        "reminder_add",
			Description: "Set a reminder",
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"title": strProp("Reminder description"),
					"when":  strProp("When to be reminded (e.g. '2025-12-31 09:00', 'tomorrow morning')"),
				},
				Required: []string{"title"},
			},
		},
		{
			Name:        "reminder_list",
			Description: "List reminders",
			InputSchema: &jsonschema.Schema{
				Type: "object",
				Properties: map[string]*jsonschema.Schema{
					"show_dismissed": boolProp("Include dismissed reminders (default: false)"),
				},
			},
		},
		{
			Name:        "reminder_dismiss",
			Description: "Dismiss a reminder",
			InputSchema: &jsonschema.Schema{
				Type:       "object",
				Properties: map[string]*jsonschema.Schema{"id": strProp("Reminder ID")},
				Required:   []string{"id"},
			},
		},
	}
}

// ---- Main ------------------------------------------------------------------

func main() {
	dataPath := os.Getenv("PERSONAL_MCP_DATA")
	if dataPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to get home dir: %v\n", err)
			os.Exit(1)
		}
		dataPath = filepath.Join(home, ".personal-mcp", "data.json")
	}

	db, err := newDB(dataPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load store from %s: %v\n", dataPath, err)
		os.Exit(1)
	}

	srv := &Server{db: db}

	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "personal-life",
		Title:   "Personal Life Manager",
		Version: "1.0.0",
	}, &mcp.ServerOptions{
		Instructions: "A personal life management assistant. Use these tools to manage your tasks, notes, journal, shopping list, and reminders.",
	})

	handlers := map[string]mcp.ToolHandler{
		"task_add":        srv.taskAdd,
		"task_list":       srv.taskList,
		"task_complete":   srv.taskComplete,
		"task_delete":     srv.taskDelete,
		"note_add":        srv.noteAdd,
		"note_list":       srv.noteList,
		"note_get":        srv.noteGet,
		"note_search":     srv.noteSearch,
		"note_delete":     srv.noteDelete,
		"journal_add":     srv.journalAdd,
		"journal_list":    srv.journalList,
		"shopping_add":    srv.shoppingAdd,
		"shopping_list":   srv.shoppingList,
		"shopping_check":  srv.shoppingCheck,
		"shopping_clear":  srv.shoppingClear,
		"reminder_add":    srv.reminderAdd,
		"reminder_list":   srv.reminderList,
		"reminder_dismiss": srv.reminderDismiss,
	}

	for _, tool := range buildTools(srv) {
		mcpServer.AddTool(tool, handlers[tool.Name])
	}

	transport := mcp.NewStdioTransport()
	if err := mcpServer.Run(context.Background(), transport); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
