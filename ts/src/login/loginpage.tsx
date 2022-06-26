import React, { useContext, useEffect } from 'react'
import { Input, Flex, Text, HStack, VStack, InputRightElement, Button, InputGroup } from '@chakra-ui/react'
import { API_URL_ROOT } from '../config/config';
import { useNavigate, Navigate } from "react-router-dom";
import { AccountInfoContext } from '../infoStore/accountInfoStore';
import { login } from '../api/apis';
import { observer } from 'mobx-react-lite';

const LoginPage = observer((prop: any) => {
    const [account, setAccount] = React.useState('');
    const [password, setPassword] = React.useState('');
    const [show, setShow] = React.useState(false);
    const accountInfo = useContext(AccountInfoContext);
    const navigate = useNavigate();
    const handleClick = () => setShow(!show);
    const handleLogin = () => {
        login(account, password).then(response => {
            if (Math.floor(response.status / 100) == 2) {
                accountInfo.setAccountName(account);
                navigate("/home", { replace: true });
            }
        });
    };

    if (accountInfo.permission == 0 || accountInfo.permission == 1 || accountInfo.permission == 2) {
        return <Navigate to="/home" replace={true} />;
    }

    return (
        <Flex justify="center">
            <VStack {...prop} justify="right">
                <VStack alignItems="flex-start" width="100%" >
                    <Text fontSize='24px'>Account</Text>
                    <Input onChange={(event) => setAccount(event.target.value)} />
                </VStack>
                <VStack alignItems="flex-start" width="100%">
                    <Text fontSize='24px'>Password</Text>
                    <InputGroup size='md'>
                        <Input type={show ? 'text' : 'password'} onChange={(event) => setPassword(event.target.value)} />
                        <InputRightElement width='4.5rem'>
                            <Button h='1.75rem' size='sm' onClick={handleClick}>
                                {show ? 'Hide' : 'Show'}
                            </Button>
                        </InputRightElement>
                    </InputGroup>
                </VStack>
                <Button bgColor="green.400" size='md' onClick={handleLogin} alignSelf="flex-end">
                    Login
                </Button>
            </VStack>
        </Flex>);
});

export default LoginPage