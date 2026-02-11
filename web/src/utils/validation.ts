import { parse } from "tldts"

const DOMAIN_REGEX = /^[a-zA-Z0-9]([a-zA-Z0-9-]*\.)+[a-zA-Z]{2,}$/

const isValidDomain = (domain: string): boolean => {
    return DOMAIN_REGEX.test(domain)
}

const isValidExtension = (domain: string): boolean => {
    const result = parse(domain)
    return result.isIcann === true
}

export { isValidDomain, isValidExtension };