package com.example.mockusers.users;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/users")
public class UserController {

    private final UserStore store;

    public UserController(UserStore store) {
        this.store = store;
    }

    @GetMapping
    public List<User> getAll() {
        return store.findAll();
    }

    @GetMapping("/{id}")
    public ResponseEntity<User> getOne(@PathVariable String id) {
        return store.findById(id)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }

    @PostMapping
    public ResponseEntity<User> create(@RequestBody User user) {
        return ResponseEntity.status(201).body(store.save(user));
    }

    @PutMapping("/{id}")
    public ResponseEntity<User> replace(@PathVariable String id, @RequestBody User user) {
        if (store.findById(id).isEmpty()) return ResponseEntity.notFound().build();
        return ResponseEntity.ok(store.save(user));
    }

    @PatchMapping("/{id}")
    public ResponseEntity<User> patch(@PathVariable String id, @RequestBody Map<String, Object> fields) {
        return store.patch(id, fields)
                .map(ResponseEntity::ok)
                .orElse(ResponseEntity.notFound().build());
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> delete(@PathVariable String id) {
        return store.delete(id)
                ? ResponseEntity.ok().<Void>build()
                : ResponseEntity.notFound().build();
    }
}
