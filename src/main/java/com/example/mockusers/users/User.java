package com.example.mockusers.users;

import com.fasterxml.jackson.annotation.JsonInclude;

public record User(
    String id,
    String type,
    String address,
    String city,
    String card_id,
    @JsonInclude(JsonInclude.Include.NON_NULL) String email
) {}
