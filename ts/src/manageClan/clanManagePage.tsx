import {
    Table,
    Thead,
    Tbody,
    Tfoot,
    Tr,
    Th,
    Td,
    TableCaption,
    TableContainer,
    Center,
    HStack,
    Button,
    AlertDialog,
    AlertDialogBody,
    AlertDialogFooter,
    AlertDialogHeader,
    AlertDialogContent,
    AlertDialogOverlay,
    useDisclosure
} from '@chakra-ui/react';
import { useEffect, useState, useRef, MouseEventHandler } from 'react';
import { fetchClanMembers, AccountInfo, promoteMember, kickMember } from '../api/apis';
const ClanManagePage = () => {
    const [members, setMembers] = useState<Array<AccountInfo>>([]);
    const { isOpen, onOpen, onClose } = useDisclosure();
    const cancelRef = useRef<any>();
    const [alertInfo, setAlertInfo] = useState<{ desc: string, display: string, onClick: MouseEventHandler<HTMLButtonElement> }>({ desc: "", display: "", onClick: () => { } });

    const init = async () => {
        let ret = await fetchClanMembers();
        if (ret) {
            ret = ret.sort((a, b) => a.Permission - b.Permission);
            setMembers(ret);
        }
    };

    useEffect(() => {
        init();
    }, []);

    const MemberRow = (info: AccountInfo) => {
        let rank = "member";
        let actions = null;
        switch (info.Permission) {
            case 1:
                rank = "leader";
                break;
            default:
                actions = (<HStack>
                    <Button backgroundColor="green.500" size='sm'
                        onClick={() => {
                            setAlertInfo({
                                desc: `Do you really want to promote ${info.Name} to leader? You cannot undo this operation.`,
                                display: "Promote",
                                onClick: () => { promoteMember(info.Name); onClose(); init(); }
                            });
                            onOpen();
                        }}>
                        Promote
                    </Button>
                    <Button backgroundColor="red.500" size='sm' onClick={() => {
                        setAlertInfo({
                            desc: `Do you really want to kick ${info.Name} from your clan? You cannot undo this operation.`,
                            display: "Kick",
                            onClick: () => { kickMember(info.Name); onClose(); init(); }
                        });
                        onOpen();
                    }}>
                        Kick
                    </Button>
                </HStack>);

        }

        return (<Tr>
            <Td>{info.Name}</Td>
            <Td>{rank}</Td>
            <Td>
                {actions}
            </Td>
        </Tr>)
    };

    return (
        <Center>
            <TableContainer width={[
                '100%', // 0-30em
                '70%', // 30em-48em
                '60%', // 48em-62em
            ]}>
                <Table variant='simple'>
                    <Thead>
                        <Tr>
                            <Th>Name</Th>
                            <Th>Rank</Th>
                            <Th>Action</Th>
                        </Tr>
                    </Thead>
                    <Tbody>
                        {members.map(info => <MemberRow {...info} key={info.Name} />)}
                    </Tbody>
                </Table>
            </TableContainer>
            <AlertDialog
                isOpen={isOpen}
                leastDestructiveRef={cancelRef}
                onClose={onClose}
            >
                <AlertDialogOverlay>
                    <AlertDialogContent>
                        <AlertDialogHeader fontSize='lg' fontWeight='bold'>
                            {alertInfo.display}
                        </AlertDialogHeader>

                        <AlertDialogBody>
                            {alertInfo.desc}
                        </AlertDialogBody>

                        <AlertDialogFooter>
                            <Button ref={cancelRef} onClick={onClose}>
                                Cancel
                            </Button>
                            <Button colorScheme='red' onClick={alertInfo.onClick} ml={3}>
                                Comfirm
                            </Button>
                        </AlertDialogFooter>
                    </AlertDialogContent>
                </AlertDialogOverlay>
            </AlertDialog>
        </Center>
    )
};

export default ClanManagePage;