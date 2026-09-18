package nsu.security.demoapplication.controller;

import nsu.security.demoapplication.model.Secret;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

import nsu.security.demoapplication.repository.SecretsJdbcRepository;
import org.springframework.web.bind.annotation.RequestHeader;
import org.springframework.http.HttpStatus;

import java.util.List;

@RestController
public class SecretsController {
    private final SecretsJdbcRepository secretsJdbcRepository;

    public SecretsController(SecretsJdbcRepository secretsJdbcRepository) {
        this.secretsJdbcRepository = secretsJdbcRepository;
    }

    @GetMapping("/secrets")
    public ResponseEntity<List<Secret>> getOrder(@RequestHeader(value = "User-Agent", required = true) String userAgent) {

        /*if (userAgent == null || !userAgent.equals("SAFE-APP")) {
            return ResponseEntity.status(HttpStatus.FORBIDDEN).build();
        }
        */
        List<Secret> secrets = secretsJdbcRepository.getSecrets();
        if (secrets.isEmpty()) {
            return ResponseEntity.notFound().build();
        }
        return ResponseEntity.ok(secrets);
    }
}
