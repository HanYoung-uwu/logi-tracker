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
    Button
} from '@chakra-ui/react';
import { useEffect, useState } from 'react';
import { fetchClanMembers, AccountInfo, promoteMember, kickMember } from '../api/apis';
const ClanManagePage = () => {
    const [members, setMembers] = useState<Array<AccountInfo>>([]);
    useEffect(() => {
        const init = async () => {
            let ret = await fetchClanMembers();
            if (ret) {
                setMembers(ret);
            }
        };
        init();
    }, []);

    const MemberRow = (info: AccountInfo) => {
        let rank = "member";
        switch (info.Permission) {
            case 1:
                rank = "leader";
                break;
        }
        return (<Tr>
            <Td>{info.Name}</Td>
            <Td>{rank}</Td>
            <Td>
                <HStack>
                    <Button backgroundColor="green.500" size='sm' onClick={() => promoteMember(info.Name)}>
                        Promote
                    </Button>
                    <Button backgroundColor="red.500" size='sm' onClick={() => kickMember(info.Name)}>
                        Kick
                    </Button>
                </HStack>
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
        </Center>
    )
};

export default ClanManagePage;