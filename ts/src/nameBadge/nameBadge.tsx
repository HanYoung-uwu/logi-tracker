import { Stack, HStack, VStack, Button, Spacer, Center, Icon, Text } from '@chakra-ui/react'
import { useEffect, useState, useContext } from 'react';
import { observer } from 'mobx-react-lite';
import { useNavigate } from 'react-router-dom'
import { MdAccountCircle } from 'react-icons/md';
import { fetchAccountInfo } from '../api/apis';
import { AccountInfoContext } from '../login/accountInfoStore';

const NameBadge = observer((props: any) => {
    let navigate = useNavigate();
    const [handleClick, setHandleClick] = useState<() => void>(() => navigate("/login", { replace: true }));
    const accountInfo = useContext(AccountInfoContext);

    useEffect(() => {
        const init = async () => {
            let info = await fetchAccountInfo();
            if (info) {
                setHandleClick(() => { console.log("TODO: ACCOUNT SETTING PAGE") });
                accountInfo.setClan(info.Clan);
                accountInfo.setPermission(info.Permission);
            } else {
                setHandleClick(() => navigate("/login", { replace: true }));
            }
        };
        if (accountInfo.name != '') {
            init();
        }
    }, [accountInfo.name]);
    return (
        <Button onClick={handleClick}>
            <HStack>
                <Icon as={MdAccountCircle} boxSize={6} />
                <Text>{accountInfo.getAccountName() == '' ? "SIGN IN" : accountInfo.getAccountName()}</Text>
            </HStack>
        </Button>
    )
});

export { NameBadge }