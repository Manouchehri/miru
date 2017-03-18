# Deployment Advice

When choosing to deploy Miru to the Internet, consider the following tips.

1. Set up a TLS Certificate to secure client-to-server communication.
   * Consider using [Let's Encrypt](https://letsencrypt.org/getting-started/) for a free and easy-to-use solution.
2. Configure your server to always use a secure HTTPS connection.
   * Use [HTTP Strict Transport Security](https://www.owasp.org/index.php/HTTP_Strict_Transport_Security_Cheat_Sheet).
     * Instructions for [nginx](https://www.nginx.com/blog/http-strict-transport-security-hsts-and-nginx/).
3. Manage your server with an account with the lowest privileges necessary.