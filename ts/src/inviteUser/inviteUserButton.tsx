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
} from '@chakra-ui/react';
import { observer } from 'mobx-react-lite';
import { useEffect, useState, useContext } from 'react';
import { AiOutlinePlus } from 'react-icons/ai';
import { fetchClanInviteLink, fetchClanAdminInviteLink } from '../api/apis';
import { WEBSITE_ROOT } from '../config/config';
import { AccountInfoContext } from '../infoStore/accountInfoStore';

const InviteUserButton = observer((props: any) => {
    const { isOpen, onOpen, onClose } = useDisclosure();
    const [inviteLink, setInviteLink] = useState('');
    const [input, setInput] = useState<HTMLInputElement | null>();
    const accountInfo = useContext(AccountInfoContext);

    useEffect(() => {
        if (accountInfo.permission == 0 || accountInfo.permission == 1) {
            let fetchLinkFunc: any;
            if (accountInfo.permission == 0) {
                fetchLinkFunc = fetchClanAdminInviteLink;
            } else if (accountInfo.permission == 1) {
                fetchLinkFunc = fetchClanInviteLink;
            }
            if (isOpen) {
                fetchLinkFunc().then((token: string) => {
                    setInviteLink(`${WEBSITE_ROOT}/invite?link=${token}`);
                });
            }
        }
    }, [isOpen]);

    if (accountInfo.permission != 0 && accountInfo.permission != 1) {
        return <></>;
    }

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
                            <InputRightElement width="4.5rem" children={<Button onClick={() => { if (input) input.select(); document.execCommand('copy'); }}>Copy</Button>} />
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
});

export default InviteUserButton;