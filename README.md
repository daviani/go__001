# GoSentry — Security Audit API

API REST d'audit de surface d'attaque externe, écrite en Go sans framework.

Analyse un domaine sur 5 axes : DNS, certificats SSL/TLS, headers de sécurité, sous-domaines et fichiers sensibles exposés.

## Stack

- **Langage** : Go (bibliothèques standard uniquement pour les scanners)
- **API** : `net/http` (serveur natif, pas de framework)
- **Documentation** : Swagger UI via [swaggo/swag](https://github.com/swaggo/swag)
- **CI/CD** : GitHub Actions (tests + build + Docker)
- **Conteneurisation** : Docker (multi-stage build)

## Scanners

| Scanner | Description | Packages Go |
|---------|-------------|-------------|
| DNS | Records A/AAAA, MX, NS, TXT | `net` |
| SSL/TLS | Certificat, émetteur, expiration | `crypto/tls` |
| Headers | HSTS, CSP, X-Frame-Options, X-Content-Type-Options | `net/http` |
| Sous-domaines | Énumération via Certificate Transparency (crt.sh) | `net/http`, `encoding/json` |
| Fichiers sensibles | Détection .env, .git/config, wp-config.php... | `net/http` |

## Démarrage rapide

```bash
# Cloner le projet
git clone https://github.com/daviani/go__001.git
cd go__001

# Configurer l'environnement
cp .env.example .env

# Lancer le serveur
go run main.go

# Le serveur écoute sur http://localhost:8082
```

## Configuration

| Variable | Description | Défaut |
|----------|-------------|--------|
| `PORT` | Port d'écoute du serveur | `8082` |

## API

| Verbe | Route | Description |
|-------|-------|-------------|
| `GET` | `/health` | Status du serveur |
| `GET` | `/scan/dns?domain=xxx` | Scan DNS |
| `GET` | `/scan/ssl?domain=xxx` | Scan certificat SSL/TLS |
| `GET` | `/scan/header?domain=xxx` | Scan headers de sécurité |
| `GET` | `/scan/subdomain?domain=xxx` | Énumération sous-domaines |
| `GET` | `/scan/sensitive?domain=xxx` | Détection fichiers sensibles |
| `GET` | `/scan/all?domain=xxx` | Lance les 5 scanners en parallèle |

Documentation interactive : [http://localhost:8082/swagger/index.html](http://localhost:8082/swagger/index.html)

## Architecture

```
go__001/
├── main.go                         # Point d'entrée
├── Dockerfile                      # Multi-stage build
├── .github/workflows/ci.yml        # CI GitHub Actions
├── docs/                           # Spec Swagger générée
└── internal/
    ├── api/
    │   ├── models.go               # Structs (Server, ScanResult...)
    │   ├── handlers.go             # Handlers HTTP + annotations Swagger
    │   ├── middleware.go           # CORS middleware
    │   └── server.go               # Routeur + démarrage serveur
    └── scanner/
        ├── scanner.go              # Interface Scanner
        ├── dns.go                  # Scanner DNS
        ├── ssl.go                  # Scanner SSL/TLS
        ├── header.go               # Scanner Headers HTTP
        ├── subdomain.go            # Scanner sous-domaines
        └── sensitive.go            # Scanner fichiers sensibles
```

## Docker

```bash
docker build -t gosentry .
docker run -p 8082:8082 gosentry
```

## Tests

```bash
go test ./... -v
```
