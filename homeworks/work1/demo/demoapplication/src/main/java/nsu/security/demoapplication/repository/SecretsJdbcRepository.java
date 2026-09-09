package nsu.security.demoapplication.repository;

import nsu.security.demoapplication.model.Secret;
import org.springframework.stereotype.Repository;

import java.sql.Connection;
import java.sql.DriverManager;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.sql.Statement;
import java.util.ArrayList;
import java.util.List;

@Repository
public class SecretsJdbcRepository {

    private static final String URL = "jdbc:postgresql://db:5432/demodb";
    private static final String USER = "demo";
    private static final String PASSWORD = "demo";

    public List<Secret> getSecrets() {

        Connection connection = null;
        Statement statement = null;
        PreparedStatement preparedStatement = null;
        ResultSet results = null;

        List<Secret> secrets = new ArrayList<>();

        try {
            connection = DriverManager.getConnection(URL, USER, PASSWORD);

            String sql = "SELECT * FROM secrets";
            preparedStatement = connection.prepareStatement(sql);
            results = preparedStatement.executeQuery();

            while (results.next()) {
                Long scretId = results.getLong("id");
                String name = results.getString("name");
                String password = results.getString("password");

                Secret secret = new Secret(name, password);
                secret.setId(scretId);
                secrets.add(secret);
            }

            return secrets;

        } catch (SQLException e) {
            throw new RuntimeException("Не удалось получить заказы через JDBC", e);

        } finally {
            closeQuietly(results);
            closeQuietly(statement);
            closeQuietly(preparedStatement);
            closeQuietly(connection);
        }
    }

    private void closeQuietly(AutoCloseable closeable) {
        if (closeable != null) {
            try {
                closeable.close();
            } catch (Exception e) {
                // намеренно игнорируем ошибку закрытия ресурса
            }
        }
    }
}