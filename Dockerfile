# ── Stage 1 : Build ──────────────────────
# Image Go complète pour compiler le code source
# "AS builder" nomme ce stage pour qu'on puisse copier depuis lui
FROM golang:1.25-alpine AS builder

# Répertoire de travail dans le conteneur (chemin absolu obligatoire)
WORKDIR /app

# Copie go.mod et go.sum EN PREMIER → optimisation cache Docker
# Si ces fichiers n'ont pas changé, Docker réutilise le cache des dépendances
COPY go.mod go.sum ./

# Télécharge les dépendances (équivalent npm install)
RUN go mod download

# Copie tout le code source (après le download pour profiter du cache)
COPY . .

# Compile le binaire Go → produit un exécutable "/app/scanner"
# -o scanner = nom du fichier de sortie (comme "npm run build" mais en binaire natif)
RUN go build -o scanner .

# ── Stage 2 : Run ────────────────────────
# Image minimale Alpine (~5MB) — pas besoin de Go pour exécuter un binaire compilé
# Le stage 1 (builder) est jeté → image finale ~17MB au lieu de ~800MB
FROM alpine:latest

# Copie uniquement le binaire compilé depuis le stage builder
# --from=builder référence le stage nommé "builder" au-dessus
COPY --from=builder /app/scanner /scanner

# Commande exécutée au lancement du conteneur
CMD ["/scanner"]
