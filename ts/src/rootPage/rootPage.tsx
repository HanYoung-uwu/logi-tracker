import {
    ChakraProvider,
    Box,
    Text,
    Link,
    VStack,
    Heading,
    Code,
    HStack,
    Center,
    Grid,
    Flex
} from "@chakra-ui/react";

const RootPage = () => {

    return (
        <Flex justify="center">
            <Center width={[
                '100%', // 0-30em
                '50%', // 30em-48em
                '25%', // 48em-62em
            ]}>
                <VStack>
                    <Heading>Welcome to Logi Tracker</Heading>
                    <Text>Currently, Logi Tracker is invitation only.</Text>
                    <Text>If you are clan leader, you can ping</Text>
                    <Text fontWeight="bold">KT_Linux</Text>
                    <Text>in </Text>
                    <Link fontSize="24px" href="https://discord.gg/jH8hC79rNX">ASEAN discord server</Link>
                    <br />
                    <Text>If your clan is already invited to Logi Tracker</Text>
                    <Text>ask your clan leader to send you an invitation</Text>
                </VStack>
            </Center>
        </Flex>)
};

export default RootPage;