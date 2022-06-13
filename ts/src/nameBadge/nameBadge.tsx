import { Stack, HStack, VStack, Button, Spacer, Center, Icon, Text } from '@chakra-ui/react'
import { useEffect, useState, useContext } from 'react';
import { observer } from 'mobx-react-lite';
import { useNavigate } from 'react-router-dom'
import { MdAccountCircle } from 'react-icons/md';
import { fetchAccountInfo } from '../api/apis';
import { AccountInfoContext } from '../login/accountInfoStore';

const NameBadge = observer((props: any) => {
    let navigate = useNavigate();
    const [handleClick, setHandleClick] = useState<any>();
    const accountInfo = useContext(AccountInfoContext);

    useEffect(() => {
        const init = async () => {
            let info = await fetchAccountInfo();
            if (info) {
                setHandleClick(() => { console.log("TODO: ACCOUNT SETTING PAGE") });
                accountInfo.name = info.Name;
            } else {
                setHandleClick(() => navigate("/login", { replace: true }));
            }
        };
        init();
    }, []);
    return (
        <Button onClick={handleClick}>
            <HStack>
                <Icon as={MdAccountCircle} boxSize={6} />
                <Text>{accountInfo.name == '' ? "SIGN IN" : accountInfo.name}</Text>
            </HStack>
        </Button>
    )
});

export { NameBadge }