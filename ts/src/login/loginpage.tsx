import React from 'react'
import { Input, Flex, Text, HStack, VStack, InputRightElement, Button, InputGroup } from '@chakra-ui/react'
import { API_URL_ROOT } from '../config/config';
import { useNavigate } from "react-router-dom";

const LoginPage = (prop: any) => {
    const [account, setAccount] = React.useState('');
    const [password, setPassword] = React.useState('');
    const [show, setShow] = React.useState(false);
    let navigate = useNavigate();
    const handleClick = () => setShow(!show);
    const handleLogin = () => {
        let url = API_URL_ROOT + "/login"
        let headers = {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        };

        const login = async () => {
            let res = await fetch(url, {
                method: "POST",
                headers: headers,
                body: JSON.stringify({
                    Name: account,
                    Password: password
                })
            });
            if (Math.floor(res.status / 100) == 2) {
                navigate("/home", { replace: true });
            }
        };
        login();
    };
    return (
        <Flex justify="center">
            <VStack {...prop} justify="right">
                <HStack>
                    <Text fontSize='24px'>Account</Text>
                    <Input onChange={(event) => setAccount(event.target.value)} />
                </HStack>
                <HStack>
                    <Text fontSize='24px'>Password</Text>
                    <InputGroup size='md'>
                        <Input type={show ? 'text' : 'password'} onChange={(event) => setPassword(event.target.value)} />
                        <InputRightElement width='4.5rem'>
                            <Button h='1.75rem' size='sm' onClick={handleClick}>
                                {show ? 'Hide' : 'Show'}
                            </Button>
                        </InputRightElement>
                    </InputGroup>
                </HStack>
                <Button h='1.75rem' size='sm' onClick={handleLogin}>
                    Login
                </Button>
            </VStack>
        </Flex>);
}

export default LoginPage