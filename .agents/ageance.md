

## 1. Philosophy

This project follows **extreme discipline in code structure, readability, and maintainability**.

Core principles:

* Single Responsibility everywhere
* High DRY enforcement
* Strict modularity
* Predictable naming
* Fail-fast, fully wrapped errors
* Observable and measurable system behavior

---

## 2. File Constraints

* **Maximum file size: 50 lines (hard limit)**
* If a file exceeds 50 lines:

  * It must be split immediately
* Files should represent **one logical unit only**

### File Design Rules

* One struct per file (preferred)
* One responsibility per file
* No mixed concerns

---

## 3. Function Design

* Each function must:

  * Do **exactly one thing**
  * Be **≤ 15 lines** (soft limit)
* Avoid nested logic > 2 levels
* No side-effects unless explicitly required

### Anti-patterns (STRICTLY NOT ALLOWED)

* God functions
* Hidden state mutations
* Mixed responsibilities
* Avoid adding any comments , code should be self-explanatory.

---

## 4. Naming Conventions (MANDATORY)

All identifiers must follow explicit suffix-based naming:

### Structs

```
user_struct
order_service_struct
```

### Constructors

```
user_const
order_service_const
```

### Methods

```
user_get_name_method
order_process_method
```

### Interfaces

```
user_repository_interface
payment_gateway_interface
```

### Variables

* Clear, explicit, no abbreviations

```
user_id
order_total_amount
```

---

## 5. Error Handling (ZERO TOLERANCE POLICY)

* Every error must be:

  * Checked
  * Wrapped
  * Contextualized

### Rules:

* NEVER ignore errors
* NEVER return raw errors
* ALWAYS wrap using:

```
fmt.Errorf("context: %w", err)
```

### Example:

```
if err != nil {
    return fmt.Errorf("failed to fetch user: %w", err)
}
```

---

## 6. DRY Enforcement

* No duplication tolerated
* If logic repeats **twice → extract**
* Shared logic must go into:

  * utilities
  * helpers
  * services

---

## 7. SOLID Principles (STRICT)

### S — Single Responsibility

Every unit must do one thing only

### O — Open/Closed

Extend via composition, not modification

### L — Liskov Substitution

Interfaces must be clean and minimal

### I — Interface Segregation

Small, focused interfaces only

### D — Dependency Inversion

Depend on abstractions, not implementations

---

## 8. Project Structure

Example:

```
/internal
    /user
        user_struct.go
        user_const.go
        user_get_method.go
    /order
    /shared
        error_utils.go
        logger.go
```

---

## 9. Observability (MANDATORY)

### Logging

* Structured logs only
* No `fmt.Println` in production code

### Metrics

* Track:

  * latency
  * error rate
  * throughput

### Monitoring

* Must include:

  * performance tracking
  * memory usage tracking
  * goroutine leak detection

---

## 10. Performance Rules

* No unnecessary allocations
* Avoid global state
* Benchmark critical paths
* Use profiling tools regularly

---

## 11. Testing (STRICT GATE)

* No code is merged without tests

### Required:

* Unit tests for all business logic
* Edge case coverage mandatory

### Commands:

```
make test
```

---

## 12. CI/CD Enforcement

* CI pipeline must:

  * Run tests
  * Run lint checks
  * Run build

* **Code cannot be merged if:**

  * Tests fail
  * Lint fails
  * Build fails

---

## 13. Makefile Requirements

Must include:

```
make build
make test
make lint
make run
```

Optional:

```
make benchmark
make profile
```

---

## 14. Code Review Rules

Every PR must ensure:

* File size limits respected
* Naming conventions followed
* Errors properly wrapped
* No duplication
* SOLID adherence

---

## 15. Non-Negotiables

* No shortcuts
* No “temporary” hacks
* No unhandled errors
* No oversized files
* No unclear naming

---

## 16. Guiding Principle

> If the code feels even slightly hard to read, it is wrong.

