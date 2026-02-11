import {Button, Flex, Input, NativeSelect, Text} from "@chakra-ui/react";
import {useState} from "react";

function ScanForm({ onScan }: { onScan: (domain: string, scanType: string) => void }) {
    const [domain, setDomain] = useState("")
    const [scanType, setScanType] = useState("all")
    return (
        <>
            <Flex p="8" maxW="600px" mx="auto" direction="column" gap="4" >
                <Text p="8" textAlign="center" >Domaine à analyser</Text>
                <Input  bg="bg.card"
                        borderColor="nord.polar3"
                        color="text.main"
                        value={domain}
                        onChange={(e) => setDomain(e.target.value)}
                        placeholder="daviani.dev"
                />
                <NativeSelect.Root>
                    <NativeSelect.Field
                        value={scanType}
                        onChange={(e) => setScanType(e.target.value)}
                        borderColor="nord.polar3" color="text.main"
                    >
                        <option value="all">Tous les scanners</option>
                        <option value="dns">DNS</option>
                        <option value="ssl">SSL/TLS</option>
                        <option value="header">Headers HTTP</option>
                        <option value="subdomain">Sous-domaines</option>
                        <option value="sensitive">Fichiers sensibles</option>
                    </NativeSelect.Field>
                </NativeSelect.Root>
                <Button bg="accent" color="nord.polar0" onClick={():any => onScan(domain, scanType)}>Scanner</Button>
            </Flex>
        </>
    )
}

export default ScanForm;