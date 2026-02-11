import {Card, Flex, Spinner} from "@chakra-ui/react";
import type {ScanResult} from "../services/scanner.ts";


function ScanResults({ results, loading }: { results: ScanResult[], loading: boolean }) {
    return (
        <Flex direction="column" gap="4" maxW="600px" mx="auto" p="8" >
            {loading ?
                <Spinner mx="auto"  display="block" color="accent" />
                :
                <Flex direction="column" gap="4">
                    {results.map((r) => (
                        <Card.Root key={r.scanner} bg="bg.card" borderColor="nord.polar3" >
                            <Card.Header textAlign="center">{r.scanner} :</Card.Header>
                            <Card.Body>{r.result}</Card.Body>
                        </Card.Root>
                    ))}
                </Flex>
            }
        </Flex>
    )
}

export default ScanResults;