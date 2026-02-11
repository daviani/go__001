import { describe, it, expect } from "vitest"
import { isValidDomain, isValidExtension } from "./validation.ts"


describe("isValidDomain", () => {
    it('accepte un domaine valide', () => {
        expect(isValidDomain("daviani.dev")).toBe(true)
    });

    it('refuse une chaine vide', () => {
        expect(isValidDomain('')).toBe(false)
    });

    it("refuse un domaine sans point", () => {
        expect(isValidDomain("pasdedot")).toBe(false)
    })

    it("refuse des caractères spéciaux", () => {
        expect(isValidDomain("!!!bizarre")).toBe(false)
    })

    it("accepte un sous-domaine", () => {
        expect(isValidDomain("www.google.com")).toBe(true)
    })
})

describe("isValidExtension", () => {
    it('accepte une extension ICANN', () => {
        expect(isValidExtension("example.fr")).toBe(true)
    })

    it('refuse une extension inventée', () => {
        expect(isValidExtension("example.fdsfd")).toBe(false)
    })

    it('accepte .com', () => {
        expect(isValidExtension("google.com")).toBe(true)
    })
})