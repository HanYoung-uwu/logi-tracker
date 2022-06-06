import { Stack, HStack, VStack, Button, Spacer, Center, Icon, Text } from '@chakra-ui/react'
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom'
import { MdAccountCircle } from 'react-icons/md';
import { fetchAccountInfo } from '../api/apis';

const NameBadge = (props: any) => {
    let navigate = useNavigate();
    const [handleClick, setHandleClick] = useState<any>();
    const [display, setDisplay] = useState("SIGN IN");

    useEffect(() => {
        const init = async () => {
            let info = await fetchAccountInfo();
            if (info) {
                setHandleClick(() => { console.log("TODO: ACCOUNT SETTING PAGE") });
                setDisplay(info.Name);
            } else {
                setHandleClick(() => navigate("/login", { replace: true }));
                setDisplay("SIGN IN");
            }
        };
        init();
    }, []);
    return (
        <Button onClick={handleClick}>
            <HStack>
                <Icon as={MdAccountCircle} boxSize={6} />
                <Text>{display}</Text>
            </HStack>
        </Button>
    )
}

export { NameBadge }