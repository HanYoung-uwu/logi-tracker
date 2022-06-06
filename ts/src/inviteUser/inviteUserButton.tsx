import { Stack, HStack, VStack, Button, Spacer, Center, Icon, Text } from '@chakra-ui/react'
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom'
import { AiOutlinePlus } from 'react-icons/ai';

const InviteUserButton = (props: any) => {
    const handleClick = () => {};
    return (<Button onClick={handleClick}>
            <HStack>
                <Icon as={AiOutlinePlus} boxSize={6} />
                <Text>Invite</Text>
            </HStack>
        </Button>);
};

export default InviteUserButton;