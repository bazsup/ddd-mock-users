package com.example.mockusers.users;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.annotation.PostConstruct;
import org.springframework.stereotype.Component;
import org.springframework.util.ResourceUtils;

import java.io.File;
import java.io.InputStream;
import java.util.*;

@Component
public class UserStore {

    private final Map<String, User> users = Collections.synchronizedMap(new LinkedHashMap<>());
    private final ObjectMapper mapper;

    public UserStore(ObjectMapper mapper) {
        this.mapper = mapper;
    }

    @PostConstruct
    public void init() throws Exception {
        JsonNode root;
        File file = new File("db.json");
        if (file.exists()) {
            root = mapper.readTree(file);
        } else {
            try (InputStream is = getClass().getResourceAsStream("/db.json")) {
                root = mapper.readTree(is);
            }
        }
        List<User> initial = mapper.readerForListOf(User.class).readValue(root.get("users"));
        initial.forEach(u -> users.put(u.id(), u));
    }

    public List<User> findAll() {
        synchronized (users) {
            return new ArrayList<>(users.values());
        }
    }

    public Optional<User> findById(String id) {
        return Optional.ofNullable(users.get(id));
    }

    public User save(User user) {
        users.put(user.id(), user);
        return user;
    }

    public Optional<User> patch(String id, Map<String, Object> fields) {
        return findById(id).map(existing -> {
            Map<String, Object> merged = mapper.convertValue(existing, mapper.getTypeFactory().constructMapType(Map.class, String.class, Object.class));
            merged.putAll(fields);
            User patched = mapper.convertValue(merged, User.class);
            users.put(id, patched);
            return patched;
        });
    }

    public boolean delete(String id) {
        return users.remove(id) != null;
    }
}
