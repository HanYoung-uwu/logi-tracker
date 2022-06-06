import { Stack, HStack, VStack, Button, Spacer, Center, Icon } from '@chakra-ui/react';
import { useNavigate } from 'react-router-dom';
import { NameBadge } from '../nameBadge/nameBadge';
import InviteUserButton from '../inviteUser/inviteUserButton';

const HeaderPanel = (props: any) => {
    let navigate = useNavigate();
    return (
        <Center h='100px' color='white'>
            <HStack height="95%" maxW={[
                '100%', // 0-30em
                '80%', // 30em-48em
                '60%', // 48em-62em
            ]} flex="1">
                <Button
                    onClick={() =>
                        navigate("/home", { replace: true })
                    }
                >
                    Logi Tracker
                </Button>
                <Spacer />
                <InviteUserButton />
                <NameBadge />
            </HStack>
        </Center>
    );
}

export default HeaderPanel