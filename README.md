# VibeScript Language Specification

**Version 0.2.2 (Draft)** · File extension: `.vibe` · July 2026

---

## 1. Overview & Design Goals

VibeScript is a general-purpose scripting language optimized for being read, written, and modified by large language models, while remaining pleasant for humans.

**Design principles (in priority order):**

1. **Token-frugal, token-predictable.** Short, common keywords and familiar C-style braces/semicolons. No boilerplate. Syntax stays close to patterns LLMs have seen billions of times — familiarity lowers error rates more than raw terseness does. Brace-delimited blocks are also robust to whitespace mangling in copy-paste and patch-splicing pipelines.
2. **Context is first-class.** Prose intent, examples, and invariants attach directly to declarations via `ctx` blocks. They are part of the AST, not throwaway comments, and tooling can extract or execute them.
3. **Nothing invisible.** Errors are values that flow through visible code paths. No exceptions, no implicit type coercion, no hidden global state. An LLM reading a function sees everything that can happen in it.
4. **One obvious way.** A single canonical formatting (`vibe fmt`): K&R braces, mandatory braces on all blocks, one loop-anything construct, one error mechanism. Fewer stylistic choices means more predictable generation and cleaner diffs.
5. **Trivially implementable.** The whole language is designed so a tree-walking interpreter fits in a weekend project. No features that require a complex runtime.

---

## 2. Lexical Structure

### 2.1 Encoding & layout

- Source files are UTF-8.
- **Semicolons terminate statements.** Newlines and indentation are not significant; whitespace is free-form.
- **Braces `{ }` delimit blocks.** Braces are mandatory, even for single-statement bodies (no dangling-else ambiguity, no goto-fail bugs).
- Canonical style (enforced by `vibe fmt`): K&R — opening brace on the same line, 2-space indent.
- **Source positions are 1-based**: the first line is line 1, the first character of a line is column 1. Error values (§7) and all tooling output use this convention.

### 2.2 Comments

```vibe
# line comment — for incidental notes only
```

Meaningful documentation belongs in `ctx` blocks (§8), not comments. Comments are discarded by the parser; `ctx` is retained.

### 2.3 Identifiers & keywords

Identifiers: `[a-zA-Z_][a-zA-Z0-9_]*`, `snake_case` by convention (enforced by `vibe fmt` warnings).

Reserved keywords (complete list):

```
fn let var if elif else for in while match ret use
ctx true false none and or not try err brk skip
```

### 2.4 Literals

| Kind  | Examples                       |
| ----- | ------------------------------ |
| int   | `42`, `-7`, `1_000_000`        |
| float | `3.14`, `-0.5`, `1e9`          |
| str   | `"hello"`, `"line1\nline2"`    |
| bool  | `true`, `false`                |
| none  | `none`                         |
| list  | `[1, 2, 3]`, `[]`              |
| map   | `{name: "Ada", age: 36}`, `{}` |

**Number rules (precise):**

- Underscores may appear between digits for readability (`1_000_000`) and are ignored (stripped) — the token's value contains digits only.
- A `.` continues a number **only when immediately followed by a digit**. Thus `1.foo` lexes as `1` `.` `foo` (a parse error), and `xs[1:3]` is unambiguous.
- An exponent is a lowercase `e` immediately followed by one or more digits (`1e9`, `3.14e2`). No `E`, no `e+`/`e-` sign, no leading-dot floats (`.5` is illegal; write `0.5`).
- Any number containing a fraction or exponent is a `float`; otherwise an `int`.

**String rules (precise):**

- Plain `"..."` strings are **single-line**: a raw newline before the closing quote is a lexical error. Use `"""..."""` for multiline content. (This localizes forgotten-quote errors to one line instead of swallowing the rest of the file.)
- A backslash shields the next character; recognized escapes are `\n`, `\t`, `\\`, `\"`. Unknown escape sequences are an error.
- An unterminated string is a lexical error reported at the opening quote.

- Map literals and blocks both use `{ }`; the parser disambiguates by position (a `{` in expression position is a map, in statement/body position a block). `{}` in expression position is the empty map.
- Map keys written bare (`name:`) are string keys. Computed keys use brackets: `{[expr]: value}`.
- **String interpolation:** `"hi {name}, you are {age + 1}"`. Any expression inside `{}`. Escape a literal brace as `{{`.
- Multiline strings: triple quotes `"""..."""`, leading common indentation stripped.

---

## 3. Types & Values

VibeScript is **gradually typed**. Annotations are optional; where present, the interpreter checks them at runtime (a future compiler may check statically).

### 3.1 Base types

```
int  float  str  bool  list  map  fn  none  any
```

- `any` is the implicit type of unannotated values.
- Parameterized forms for precision when it matters: `list[int]`, `map[str]` (value type; keys are always `str` or `int`).

### 3.2 Type modifiers

| Form | Meaning                                  |
| ---- | ---------------------------------------- |
| `T?` | `T` or `none` (optional)                 |
| `T!` | `T` or an error value (fallible, see §7) |

**Rules:** `!` may appear only on function return types, outermost, at most once. `?` may wrap any base type. `-> str?!` reads as "returns (`str` or `none`), or an error."

### 3.3 No implicit coercion

- `1 + "2"` is a type error. Convert explicitly: `str(1)`, `int("2")?`, `float(x)`.
- Conditions (`if`, `while`, `and`, `or`, `not`) require `bool`. There is **no truthiness**: `if (list)` is an error, write `if (len(list) > 0)`. Parentheses around conditions are optional; `vibe fmt` omits them.
- Equality `==` is structural for all types; comparing different types is `false`, never an error.

### 3.4 Mutability

- `let` bindings are immutable (cannot be reassigned).
- `var` bindings are reassignable.
- Lists and maps are mutable in place regardless of binding kind; `let` freezes the _binding_, not the value. (Kept simple deliberately — deep immutability is out of scope for v0.2.)

---

## 4. Bindings & Expressions

```vibe
let name = "Ada";            # immutable, inferred type
var count: int = 0;          # mutable, annotated
count = count + 1;
```

### 4.1 Operators (by precedence, high → low)

```
()  []  .  ?           # call, index, field access, error propagation
-   not                # unary
*   /   //  %          # multiplicative (int / int -> float; // is int division)
+   -                  # additive; + also concatenates str and list
==  !=  <  <=  >  >=   # comparison
and or                 # logical, short-circuit
|>                     # pipeline
```

There is no C ternary `?:` — the postfix `?` is reserved for error propagation. Use the `if` expression (§6) instead.

### 4.2 Indexing & access

```vibe
xs[0];          # list index; negative indexes from end: xs[-1]
xs[1:3];        # slice (start inclusive, end exclusive)
m.name;         # map field access, sugar for m["name"]
m["name"];      # map access
xs[99];         # out-of-range access returns an error value (see §7)
m.missing;      # missing key returns none (maps are lenient; lists are strict)
```

### 4.3 Pipeline

`x |> f(y)` is sugar for `f(x, y)` — the piped value becomes the **first** argument. Pipelines encourage flat, linear dataflow that reads top-to-bottom:

```vibe
let result = load("users.json")?
  |> filter(|u| u.active)
  |> map(|u| u.email)
  |> sort();
```

### 4.4 Lambdas

```vibe
let double = |x| x * 2;
let add = |a, b| a + b;
```

Single-expression only. Multi-statement logic must be a named `fn` — this pushes complexity into named, `ctx`-documentable units.

---

## 5. Functions

```vibe
fn greet(name: str, punct: str = "!") -> str {
  "hello {name}{punct}"
}
```

- **A block evaluates to the value of its last statement.** A trailing semicolon on the last statement is optional and does not change the value (unlike Rust — semicolons are pure terminators here, never semantic). `ret expr;` returns early; bare `ret;` returns `none`.
- Default parameter values are allowed; callers may use positional or named arguments: `greet("Ada", punct: "?")`.
- Functions are values and may be passed, stored, and returned. Closures capture by reference.
- No overloading, no variadics in v0.2. (A single `list` or `map` parameter covers those cases.)

---

## 6. Control Flow

```vibe
if score >= 90 {
  grade = "A";
} elif score >= 80 {
  grade = "B";
} else {
  grade = "C";
}
```

`if` is also an expression; each branch block yields its last value:

```vibe
let label = if ok { "yes" } else { "no" };
```

Loops:

```vibe
for user in users {            # iterates list values
  send(user.email)?;
}

for i, user in users {         # with index
  print("{i}: {user.name}");
}

for key, val in config {       # map iteration, insertion order (guaranteed)
  print("{key} = {val}");
}

while count < 10 {
  count = count + 1;
}
```

- `brk;` breaks the innermost loop; `skip;` continues to the next iteration.

### 6.1 `match`

```vibe
match status {
  "active" -> enable(user);
  "banned" -> ret err("banned user");
  int n    -> retry(n);          # type pattern with binding
  _        -> { log("unknown"); notify(); }
}
```

- Arms are `pattern -> expr;` or `pattern -> { block }` (no semicolon after a block arm).
- Patterns: literal values, type patterns (`int n`, `str s`), error pattern (`err e`), and wildcard `_`. Arms are checked top to bottom; `match` is an expression whose value is the matched arm's value. No exhaustiveness requirement in v0.2 (unmatched value yields `none`).

---

## 7. Error Handling

Errors are ordinary **values**, never exceptions. This keeps every possible control path visible in the source text.

### 7.1 Creating and typing errors

```vibe
ctx "Parse a non-negative age from text" {
  @example parse_age("42") -> 42
  @example parse_age("-1") -> err
}
fn parse_age(s: str) -> int! {
  let n = int(s)?;
  if n < 0 {
    ret err("age must be non-negative, got {n}");
  }
  n
}
```

- `err(msg)` constructs an error value. It automatically carries the source file and line.
- A function that can fail declares `-> T!`. Returning an error from a function not marked `!` is a runtime error (catches drift between signature and body).

### 7.2 Propagation with `?`

`expr?` unwraps a fallible value: on success it yields the inner value; on error it **returns that error from the enclosing function** immediately.

```vibe
fn load_config(path: str) -> map! {
  let text = read_file(path)?;     # bubbles up file errors
  let cfg = parse_json(text)?;     # bubbles up parse errors
  cfg
}
```

### 7.3 Recovery with `try`

```vibe
let port = try int(env("PORT")) else 8080;      # fallback value
let cfg = try load_config("app.cfg") else {};   # any error -> default
```

`try EXPR else FALLBACK` evaluates to `EXPR`'s success value, or `FALLBACK` if it errored. To inspect the error, bind it:

```vibe
match parse_age(input) {
  err e -> print("bad input: {e.msg}");
  int n -> save(n);
}
```

An error value is a map-like value with fields `.msg`, `.file`, `.line`.

**Rule of thumb encoded in the design:** `?` to pass errors up, `try ... else` to absorb them, `match` to inspect them. There is no fourth way.

---

## 8. The Context System (`ctx`)

The signature feature of VibeScript. A `ctx` block attaches structured context to the declaration that immediately follows it (a `fn`, a top-level `let`/`var`, or a module — see §9). It is parsed into the AST and preserved by all tooling.

### 8.1 Syntax

```vibe
ctx "One-line summary of what this does" {
  @intent Free prose: why this exists, assumptions, edge cases,
          links to specs — anything a future reader (human or LLM)
          needs.
  @example dedupe([1, 1, 2]) -> [1, 2]
  @example dedupe([]) -> []
  @invariant output preserves first-occurrence order
}
fn dedupe(xs: list) -> list {
  ...
}
```

- The summary string is required. The braced tag section is optional; a bare `ctx "summary"` (no braces) is valid.
- **Inside `ctx` braces, content is line-oriented:** each line starting with `@` begins a tag; following lines without `@` continue the previous tag. No semicolons are used inside `ctx` blocks — tag prose may freely contain `;`, `{`, `}`, quotes, and `#`. The block is terminated by a `}` appearing as the **first non-whitespace character of a line** (which is where canonical formatting places it); a `}` anywhere else, such as in an `@example` map literal, is ordinary prose.
- Recognized tags: `@intent`, `@example`, `@invariant`, `@see` (free-form reference). Unknown tags are preserved but ignored by tooling.

### 8.2 Executable examples

Each `@example` has the form `CALL -> EXPECTED`, where `EXPECTED` is a literal value or the keyword `err`. The command `vibe test file.vibe` runs every example and reports mismatches. Examples are therefore simultaneously:

1. documentation for humans,
2. few-shot input/output pairs for LLMs, and
3. a regression test suite — with zero extra files.

### 8.3 Context extraction

`vibe ctx file.vibe` emits a compact **context digest**: every public signature plus its `ctx` block, with bodies omitted. This is the intended way to embed a codebase into an LLM prompt cheaply — a "header file for language models":

```
fn dedupe(xs: list) -> list
  "Remove duplicates, keep first occurrence"
  ex: dedupe([1, 1, 2]) -> [1, 2]
fn load_config(path: str) -> map!
  "Read and parse app config"
```

The digest format is stable and line-oriented so tools can filter it with grep.

---

## 9. Modules

- **One file = one module.** The module name is the filename without extension.
- `use utils;` imports `utils.vibe` from the same directory (or the interpreter's search path) and exposes its top-level names under the `utils.` prefix: `utils.dedupe(xs)`.
- `use utils: dedupe, flatten;` imports specific names unprefixed.
- Names starting with `_` are module-private and not importable.
- A file may open with a `ctx` block before any declaration; it documents the module itself and appears first in the context digest.
- Imports are evaluated once and cached; circular imports are a load-time error.

---

## 10. Built-ins (Minimal Standard Library)

All built-ins are global; no imports needed. This is the complete v0.2 surface.

**Core:** `print(x)` · `len(x)` · `type(x) -> str` · `str(x)` · `int(x) -> int!` · `float(x) -> float!` · `range(a, b) -> list`

**Lists:** `push(xs, x)` · `pop(xs) -> any!` · `sort(xs) -> list` · `rev(xs) -> list` · `map(xs, f)` · `filter(xs, f)` · `reduce(xs, f, init)` · `join(xs, sep) -> str` · `has(xs, x) -> bool`

**Maps:** `keys(m) -> list` · `vals(m) -> list` · `del(m, k)` · `has(m, k) -> bool`

**Strings:** `split(s, sep) -> list` · `trim(s)` · `lower(s)` · `upper(s)` · `replace(s, old, new)` · `has(s, sub) -> bool`

**IO / system (all fallible):** `read_file(path) -> str!` · `write_file(path, s) -> none!` · `env(name) -> str!` · `args() -> list` · `parse_json(s) -> any!` · `to_json(x) -> str`

Note the deliberate reuse: `has` works on lists, maps, and strings; `len` on all collections. Fewer names to remember (or for a model to guess wrong).

---

## 11. Tooling Conventions

These commands are part of the language contract, even before all are implemented:

| Command                | Purpose                                                                                                                                                                     |
| ---------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `vibe run file.vibe`   | Execute a script.                                                                                                                                                           |
| `vibe test file.vibe`  | Run all `@example`s as tests.                                                                                                                                               |
| `vibe fmt file.vibe`   | Rewrite to the single canonical style: K&R braces, 2-space indent, one statement per line, one blank line between top-level declarations, no parentheses around conditions. |
| `vibe ctx file.vibe`   | Emit the context digest (§8.3).                                                                                                                                             |
| `vibe check file.vibe` | Parse + validate types/annotations without running.                                                                                                                         |

Canonical formatting matters for LLMs: deterministic layout makes diffs minimal and generated code trivially mergeable. Because layout is not semantic, `vibe fmt` can always repair whitespace-mangled but syntactically valid code.

---

## 12. Grammar (EBNF, abridged)

```ebnf
program     = { ctx_block | statement } ;
ctx_block   = "ctx" string [ "{" { tag_line } "}" ] ;
tag_line    = "@" ident text_to_eol { continuation_line } ;

statement   = fn_decl | let_decl | var_decl | assign | if_stmt
            | for_stmt | while_stmt | match_expr | use_stmt
            | "ret" [ expr ] ";" | "brk" ";" | "skip" ";"
            | expr ";" ;

fn_decl     = "fn" ident "(" [ params ] ")" [ "->" type ] block ;
params      = param { "," param } ;
param       = ident [ ":" type ] [ "=" expr ] ;
type        = base_type [ "?" ] [ "!" ] | base_type "[" type "]" ;
base_type   = "int" | "float" | "str" | "bool" | "list"
            | "map" | "fn" | "none" | "any" ;

let_decl    = "let" ident [ ":" type ] "=" expr ";" ;
var_decl    = "var" ident [ ":" type ] "=" expr ";" ;
assign      = target "=" expr ";" ;

if_stmt     = "if" expr block { "elif" expr block } [ "else" block ] ;
for_stmt    = "for" ident [ "," ident ] "in" expr block ;
while_stmt  = "while" expr block ;
match_expr  = "match" expr "{" { arm } "}" ;
arm         = pattern "->" ( expr ";" | block ) ;
pattern     = literal | base_type ident | "err" ident | "_" ;
use_stmt    = "use" ident [ ":" ident { "," ident } ] ";" ;

block       = "{" { statement } "}" ;       # value = last statement's value

expr        = pipe_expr ;
pipe_expr   = or_expr { "|>" call } ;
or_expr     = and_expr { "or" and_expr } ;
and_expr    = cmp_expr { "and" cmp_expr } ;
cmp_expr    = add_expr [ cmp_op add_expr ] ;
add_expr    = mul_expr { ("+" | "-") mul_expr } ;
mul_expr    = unary { ("*" | "/" | "//" | "%") unary } ;
unary       = [ "-" | "not" ] postfix ;
postfix     = primary { call_args | index | "." ident | "?" } ;
primary     = literal | ident | lambda | "(" expr ")"
            | "if" expr block { "elif" expr block } "else" block
            | "try" expr "else" expr | "err" "(" expr ")"
            | match_expr ;
lambda      = "|" [ params ] "|" expr ;
```

Note: a `{` in expression position parses as a map literal; in statement/body position as a block.

---

## 13. Complete Example

```vibe
ctx "CLI tool: report the top email domains among active users" {
  @see users.json format: [{name, email, active}]
}

use text_stats: tally;

ctx "Extract the domain part of an email address" {
  @example domain("ada@lovelace.io") -> "lovelace.io"
  @example domain("nope") -> err
  @invariant result is lowercase
}
fn domain(email: str) -> str! {
  let parts = split(email, "@");
  if len(parts) != 2 {
    ret err("not an email: {email}");
  }
  lower(parts[1])
}

fn main() -> none! {
  let users = read_file("users.json")? |> parse_json()?;
  let domains = users
    |> filter(|u| u.active)
    |> map(|u| try domain(u.email) else "invalid")
    |> tally();                      # -> map of domain -> count

  for name, count in domains {
    print("{name}: {count}");
  }
}

main()?;
```

---

## 14. Deliberate Omissions (v0.2)

Classes/objects (maps + functions suffice) · exceptions · inheritance · macros · async/concurrency · operator overloading · a C ternary `?:` · variadics · generics beyond `list[T]`/`map[T]` · package manager. Each omission removes a dimension of stylistic variance an LLM would otherwise have to guess about. Revisit only with evidence.

---

## Changelog

- **v0.2.2** — Fixed `ctx` block termination rule: the block ends at a `}` that is the first non-whitespace character of a line, instead of forbidding `}` in prose with a `}}` escape. The old rule made `@example` lines containing map literals (`@example f() -> {a: 1}`) unwritable without escaping — exactly the kind of trap the language exists to avoid.

- **v0.2.1** — Lexical clarifications discovered during lexer implementation: 1-based source positions; precise number rules (underscores stripped, `.`/`e` require a following digit, no exponent signs or leading-dot floats); precise string rules (plain strings are single-line, escape set fixed at `\n \t \\ \"`, unterminated strings are lexical errors).

- **v0.2** — Switched from indentation-based blocks to C-style braces and semicolons. Rationale: robustness to whitespace corruption in copy-paste and patch-splicing pipelines; layout is now fully non-semantic and always repairable by `vibe fmt`. Semicolons are pure terminators (never semantic, unlike Rust). `then` keyword removed; `if` expression now uses brace form.
- **v0.1** — Initial draft (indentation-based).

_End of specification — VibeScript v0.2 (Draft)._
