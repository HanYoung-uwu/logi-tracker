import {
    Stack, HStack, VStack, Button, Spacer, Center, Icon, Text, Modal,
    ModalOverlay,
    ModalContent,
    ModalHeader,
    ModalFooter,
    ModalBody,
    ModalCloseButton,
    useDisclosure,
    Input,
    InputGroup,
    InputRightElement
} from '@chakra-ui/react'
import { useEffect, useState } from 'react';
import { AiOutlinePlus } from 'react-icons/ai';
import { fetchClanInviteLink } from '../api/apis';
import { WEBSITE_ROOT } from '../config/config';

const InviteUserButton = (props: any) => {
    const { isOpen, onOpen, onClose } = useDisclosure();
    const [inviteLink, setInviteLink] = useState('');
    const [input, setInput] = useState<HTMLInputElement | null>();

    useEffect(() => {
        if (isOpen) {
            fetchClanInviteLink().then((token) => {
                setInviteLink(`${WEBSITE_ROOT}/invite?link=${token}`);
            });
        }
    }, [isOpen]);

    return (<>
        <Button onClick={onOpen}>
            <HStack>
                <Icon as={AiOutlinePlus} boxSize={6} />
                <Text>Invite</Text>
            </HStack>
        </Button>
        <Modal isOpen={isOpen} onClose={onClose}>
            <ModalOverlay />
            <ModalContent>
                <ModalHeader>Invitation Link</ModalHeader>
                <ModalCloseButton />
                <ModalBody>
                    <VStack>
                        <Text fontSize='lg'>
                            The invitation link will expire in 24 hours.
                        </Text>
                        <InputGroup>
                            <Input ref={(input) => setInput(input)} pr='4.5rem' variant='filled' isReadOnly={true} value={inviteLink} />
                            <InputRightElement width="4.5rem" children={<Button onClick={() => {if (input) input.select(); document.execCommand('copy');}}>Copy</Button>} />
                        </InputGroup>
                    </VStack>
                </ModalBody>

                <ModalFooter>
                    <Button colorScheme='blue' mr={3} onClick={onClose}>
                        Close
                    </Button>
                </ModalFooter>
            </ModalContent>
        </Modal>
    </>);
};

export default InviteUserButton;