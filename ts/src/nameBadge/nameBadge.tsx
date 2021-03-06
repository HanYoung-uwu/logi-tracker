import {
    HStack,
    Button,
    Icon,
    Text,
    Popover,
    PopoverTrigger,
    PopoverContent,
    PopoverHeader,
    PopoverBody,
    PopoverFooter,
    PopoverArrow,
    PopoverCloseButton,
    useDisclosure
} from '@chakra-ui/react'
import { useEffect, useState, useContext } from 'react';
import { observer } from 'mobx-react-lite';
import { useNavigate } from 'react-router-dom'
import { MdAccountCircle } from 'react-icons/md';
import { fetchAccountInfo, logout } from '../api/apis';
import { AccountInfoContext } from '../infoStore/accountInfoStore';

const NameBadge = observer((props: any) => {
    let navigate = useNavigate();
    const accountInfo = useContext(AccountInfoContext);
    const { onOpen, onClose, isOpen } = useDisclosure();
    const handleClick = () => {
        if (accountInfo.getAccountName() == '') {
            navigate("/login", { replace: true });
        } else {
            onOpen();
        }
    };

    const handleLogout = () => {
        logout();
        accountInfo.setAccountName('');
        accountInfo.setClan('');
        accountInfo.setPermission(-1);
        navigate("/login", { replace: true });
        onClose();
    };

    useEffect(() => {
        const init = async () => {
            let info = await fetchAccountInfo();
            if (info) {
                accountInfo.setAccountName(info.Name);
                accountInfo.setClan(info.Clan);
                accountInfo.setPermission(info.Permission);
            } else {
                accountInfo.setPermission(-2);
            }
        };
        init();
    }, [accountInfo.name]);
    return (
        <Popover isOpen={isOpen}>
            <PopoverTrigger>
                <Button onClick={handleClick}>
                    <HStack>
                        <Icon as={MdAccountCircle} boxSize={6} />
                        <Text>{accountInfo.getAccountName() == '' ? "SIGN IN" : accountInfo.getAccountName()}</Text>
                    </HStack>
                </Button>
            </PopoverTrigger>
            <PopoverContent>
                <PopoverArrow />
                <PopoverCloseButton onClick={() => onClose()} />
                <PopoverHeader fontSize="20px" alignSelf="flex-start">{accountInfo.getAccountName()}</PopoverHeader>
                <PopoverBody>{`Clan: ${accountInfo.getClan()}`}</PopoverBody>
                <PopoverFooter display='flex' justifyContent="flex-end">
                    <Button
                        bgColor="red.500"
                        onClick={handleLogout}>Logout
                    </Button>
                </PopoverFooter>
            </PopoverContent>
        </Popover>
    )
});

export { NameBadge }