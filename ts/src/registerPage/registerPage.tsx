import { Flex, HStack, VStack, Button, Spacer, Center, Input, Text, InputGroup, InputRightElement, Alert, AlertIcon } from '@chakra-ui/react'
import React, { useEffect, useState, useContext } from 'react';
import { checkNameExist, fetchInviteAccountInfo, InviteAccountInfo } from '../api/apis';
import { useNavigate, useSearchParams } from "react-router-dom";
import { API_URL_ROOT } from '../config/config';
import { AccountInfoContext } from '../login/accountInfoStore';

const RegisterPage = (props: any) => {

    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [accountInfo, setAccountInfo] = useState<InviteAccountInfo>({
        Clan: "",
        Permission: -1
    });
    const globalAccountInfo = useContext(AccountInfoContext);
    const [show, setShow] = useState(false);
    const [params, setParams] = useSearchParams();
    const [nameValid, setNameValid] = useState(false);
    const [passwordValid, setPasswordValid] = useState(false);

    const navigate = useNavigate();

    useEffect(() => {
        let token = params.get("link");
        if (token) {
            fetchInviteAccountInfo(token).then(info => { if (info) setAccountInfo(info); });
        }
    }, []);

    const handleClick = () => {
        setShow(!show);
    };
    const handleCreate = () => {
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
                    Password: password,
                    Clan: accountInfo.Clan,
                    Token: params.get("link")
                })
            });
            if (Math.floor(res.status / 100) == 2) {
                globalAccountInfo.setAccountName(username);
                navigate("/home", { replace: true });
            }
        };
        register();
    };

    const handleUserName = (event: React.ChangeEvent<HTMLInputElement>) => {
        let name = event.target.value;
        if (name != '') {
            checkNameExist(name).then(exist => setNameValid(!exist));
        } else {
            setNameValid(false);
        }
        setUsername(name);
    };

    const handlePassword = (event: React.ChangeEvent<HTMLInputElement>) => {
        let m_password = event.target.value;
        setPassword(m_password);
        if (m_password.length < 8) {
            setPasswordValid(false);
        } else {
            setPasswordValid(true);
        }
    };

    const handleClan = (event: React.ChangeEvent<HTMLInputElement>) => {
        let info = accountInfo;
        info.Clan = event.target.value;
        setAccountInfo(info);
    };

    const getClanInput = () => {
        if (accountInfo.Permission == 4) {
            return (
                <VStack width="100%">
                    <Text fontSize='24px' alignSelf="flex-start">Clan</Text>
                    <Input onChange={handleClan} />
                </VStack>
            );
        } else {
            return <></>;
        }
    };

    return (<Center>
        <VStack width={[
            '100%', // 0-30em
            '50%', // 30em-48em
            '30%', // 48em-62em
        ]} alignItems="flex-start">
            <Text fontSize="32px">{accountInfo.Permission == 4 ? "Create your own clan!" : `You've been invited to ${accountInfo.Clan}`}</Text>
            {getClanInput()}
            <Text fontSize='24px'>Account</Text>
            <HStack width="100%">
                <Input onChange={handleUserName}
                    isInvalid={!nameValid}
                    errorBorderColor='crimson' />
                <Alert fontSize='sm' fontStyle="italic" width={nameValid ? "10%" : "80%"} status={nameValid ? "success" : "error"}>
                    <AlertIcon />
                    {nameValid ? "" : username == '' ? "name empty" : "name already existed"}
                </Alert>
            </HStack>
            <Text fontSize='24px'>Password</Text>
            <HStack width="100%">
                <InputGroup size='md'>
                    <Input type={show ? 'text' : 'password'}
                        onChange={handlePassword}
                        isInvalid={!passwordValid}
                        errorBorderColor='crimson' />
                    <InputRightElement width='4.5rem'>
                        <Button h='1.75rem'
                            size='sm'
                            onClick={handleClick}>
                            {show ? 'Hide' : 'Show'}
                        </Button>
                    </InputRightElement>
                </InputGroup>
                <Alert fontSize='sm' fontStyle="italic" width={passwordValid ? "10%" : "80%"} status={passwordValid ? "success" : "error"}>
                    <AlertIcon />
                    {passwordValid ? "" : "password must be at least 8 characters long"}
                </Alert>
            </HStack>
            <Button h='1.75rem' onClick={handleCreate}
                alignSelf="flex-end"
                disabled={!nameValid && !passwordValid}>
                Register
            </Button>
        </VStack>
    </Center>);
}

export default RegisterPage