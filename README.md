# GoSentry — Security Audit API

API REST d'audit de surface d'attaque externe, écrite en Go sans framework.

Analyse un domaine sur 5 axes : DNS, certificats SSL/TLS, headers de sécurité, sous-domaines et fichiers sensibles exposés.

## Stack

- **Langage** : Go (bibliothèques standard uniquement pour les scanners)
- **API** : `net/http` (serveur natif, pas de framework)
- **Documentation** : Swagger UI via [swaggo/swag](https://github.com/swaggo/swag)
- **Frontend** : React + TypeScript + Chakra UI (thème Nord)
- **Tests frontend** : Vitest
- **CI/CD** : GitHub Actions (tests Go + tests front + build + Docker)
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
├── .github/workflows/ci.yml        # CI GitHub Actions (Go + front)
├── docs/                           # Spec Swagger générée
├── internal/
│   ├── api/
│   │   ├── models.go               # Structs (Server, ScanResult...)
│   │   ├── handlers.go             # Handlers HTTP + annotations Swagger
│   │   ├── middleware.go           # CORS middleware
│   │   └── server.go               # Routeur + démarrage serveur
│   └── scanner/
│       ├── scanner.go              # Interface Scanner
│       ├── dns.go                  # Scanner DNS
│       ├── ssl.go                  # Scanner SSL/TLS
│       ├── header.go               # Scanner Headers HTTP
│       ├── subdomain.go            # Scanner sous-domaines
│       └── sensitive.go            # Scanner fichiers sensibles
└── web/                            # Frontend React
    ├── src/
    │   ├── App.tsx                 # Orchestrateur principal
    │   ├── components/             # Header, ScanForm, ScanResults
    │   ├── services/scanner.ts     # Client API avec URLSearchParams
    │   └── utils/validation.ts     # Validation domaine (regex + tldts)
    └── package.json
```

## Docker

```bash
docker build -t gosentry .
docker run -p 8082:8082 gosentry
```

## Tests

```bash
# Tests Go uniquement
go test ./... -v

# Tests frontend uniquement
npm run test-web

# Tous les tests (Go + front en parallèle)
npm test
```

## Développement

```bash
# Lancer le backend Go + frontend React en parallèle
npm run dev
```

## Pistes d'amélioration

### Sécurité & robustesse

- **Timeouts HTTP** : les scanners utilisent `http.Get` sans timeout — un domaine lent bloque la goroutine indéfiniment. Fix : `http.Client{Timeout: 10 * time.Second}` ou `context.WithTimeout`
- **Validation côté backend** : le handler vérifie seulement `domain == ""`. Ajouter une validation de format (regex, longueur max 253 chars) pour rejeter les inputs malformés avant de lancer les scanners
- **Recovery dans les goroutines** : si un scanner panic dans `/scan/all`, le channel bloque et la requête hang. Ajouter `defer recover()` dans chaque goroutine ou utiliser un buffered channel
- **CORS configurable** : l'origin est hardcodée à `localhost:3000`. Passer à une variable d'environnement `CORS_ORIGIN`

### Scanner fichiers sensibles

- **Élargir la liste de paths** : ajouter `/.env.local`, `/.env.production`, `/backup.sql`, `/phpinfo.php`, `/server-status`, `/.htpasswd`, `/.svn/entries`
- **Détecter les redirections** : actuellement seul le status 200 est vérifié. Un 301/302 vers un fichier sensible devrait aussi être signalé
- **Analyser le contenu** : vérifier que la réponse contient bien du contenu sensible (pas une page 200 générique type "Not Found" avec status 200)
- **Gérer les erreurs individuellement** : si un path échoue, continuer les autres au lieu de `return` sur la première erreur

### Scanners existants

- **Headers** : ajouter `X-Content-Type-Options`, `Referrer-Policy`, `Permissions-Policy` au scan
- **DNS** : rendre la gestion d'erreur cohérente (actuellement A fatal, MX/NS/TXT silencieux)
- **SSL** : message d'erreur "no client certificate" → corriger en "no peer certificate"

### Infrastructure

- **PostgreSQL** : stocker les résultats de scan pour historique et comparaison
- **Rate limiting** : protéger les endpoints contre l'abus
- **API_URL configurable** : côté front, utiliser `import.meta.env.VITE_API_URL` au lieu du port hardcodé
- **Swagger** : régénérer automatiquement la doc dans la CI pour éviter le drift
