import { Flex, HStack, VStack, Button, Spacer, Center, Input, Text, InputGroup, InputRightElement } from '@chakra-ui/react'
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom'
import { MdAccountCircle } from 'react-icons/md';
import { fetchAccountInfo } from '../api/apis';
import {useSearchParams } from "react-router-dom";
import { API_URL_ROOT } from '../config/config';

const RegisterPage = (props: any) => {

    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [show, setShow] = useState(false);
    const [params, setParams] = useSearchParams();

    const handleClick = () => {
        setShow(!show);
    };
    const handleCreate = () => {
        let token = params.get("link");
        document.cookie = `token=${token}`;
        let url = API_URL_ROOT + "/register"
        let headers = {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        };

        const register = async () => {
            let res = await fetch(url, {
                method: "POST",
                headers: headers,
                body: JSON.stringify({
                    Name: username,
                    Password: password
                })
            });
            if (Math.floor(res.status / 100) == 2) {
                console.log("created!");
            }
        };
        register();
    };

    return (<Center>
        <VStack width={[
            '100%', // 0-30em
            '50%', // 30em-48em
            '25%', // 48em-62em
        ]} alignItems="flex-start">
            <Text fontSize='24px'>Account</Text>
            <Input onChange={(event) => setUsername(event.target.value)} />
            <Text fontSize='24px'>Password</Text>
            <InputGroup size='md'>
                <Input type={show ? 'text' : 'password'} onChange={(event) => setPassword(event.target.value)} />
                <InputRightElement width='4.5rem'>
                    <Button h='1.75rem' size='sm' onClick={handleClick}>
                        {show ? 'Hide' : 'Show'}
                    </Button>
                </InputRightElement>
            </InputGroup>
            <Button h='1.75rem' onClick={handleCreate} alignSelf="flex-end">
                Register
            </Button>
        </VStack>
    </Center>);
}

export default RegisterPage