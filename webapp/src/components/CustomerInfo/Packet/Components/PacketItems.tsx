import React, {createContext, useContext} from 'react';

import {GetInputProps} from '@mantine/form/lib/types';

import {CustomerPacketValues} from '@/types/customers';

import {FormTextInput} from '@/components/form/FormTextInput';

import {FormSelect} from '@/components/form/FormDropdown';

import {FormNumberInput} from '@/components/form/FormNumberInput';

import {FormMultiSelect} from '@/components/form/FormMultiSelect';

import {parseJson} from '@/helpers/general';

import {Item} from './Item';

export const transformKeys: Record<keyof CustomerPacketValues, string> = {

    // metrics
    activeUsers: 'Active Users',
    dailyActiveUsers: 'Daily Active Users',
    monthlyActiveUsers: 'Monthly Active Users',
    inactiveUserCount: 'Inactive User Count',
    totalChannels: 'Total Channels',
    totalPosts: 'Total Posts',
    totalTeams: 'Total Teams',
    licenseSupportedUsers: 'License Supported Users',

    // environment
    deploymentType: 'Deployment Type',
    serverArch: 'Server Arch',
    serverOS: 'Server OS',
    ldapProvider: 'LDAP Provider',
    databaseSchemaVersion: 'Database Schema Version',
    databaseType: 'Database Type',
    databaseVersion: 'Database Version',
    elasticServerVersion: 'Elastic Server Version',
    fileDriver: 'File Driver',
    hostingType: 'Hosting Type',
    metricService: 'Metric Service',
    samlProvider: 'SAML Provider',
    version: 'Version',

    // general
    mobileApp: 'Mobile App',
    licensedTo: 'Licensed To',
    productsInUse: 'Products In Use',

};

export const defaultPacket: CustomerPacketValues = {
    activeUsers: 0,
    dailyActiveUsers: 0,
    monthlyActiveUsers: 0,
    inactiveUserCount: 0,
    totalChannels: 0,
    totalPosts: 0,
    totalTeams: 0,
    licenseSupportedUsers: 0,
    deploymentType: '',
    serverArch: '',
    serverOS: '',
    ldapProvider: '',
    databaseSchemaVersion: '',
    databaseType: '',
    databaseVersion: '',
    elasticServerVersion: '',
    fileDriver: '',
    hostingType: '',
    metricService: '',
    samlProvider: '',
    version: '',
    mobileApp: '',
    licensedTo: '',
    productsInUse: [],
};

type PacketContextType = {
    packet: CustomerPacketValues;
    editing: boolean;
    getInputProps: GetInputProps<CustomerPacketValues>
};

const PacketContext = createContext<PacketContextType>({
    packet: defaultPacket,
    editing: false,
    getInputProps: () => {
        return {
            // eslint-disable-next-line no-empty-function
            onChange: () => {},
            value: '',
        };
    },
});

export const PacketProvider: React.FC<{
    packet: CustomerPacketValues,
    editing: boolean,
    getInputProps: GetInputProps<CustomerPacketValues>
}> = ({children, packet, editing, getInputProps}) => {
    return (
        <PacketContext.Provider value={{packet, editing, getInputProps}}>
            {children}
        </PacketContext.Provider>
    );
};

export const PacketTextInput = ({
    packetKey,
}: {packetKey: keyof CustomerPacketValues}) => {
    const {packet, editing, getInputProps} = useContext(PacketContext);
    const title = transformKeys[packetKey];
    const value = String(packet[packetKey]);
    return (
        <Item
            title={title}
            value={value}
            type='string'
            editing={editing}
            editComponent={
                <FormTextInput
                    formKey={packetKey}
                    getInputProps={getInputProps}
                    key={packetKey}
                    placeholder={title}
                />
            }
        />
    );
};
export const PacketNumberField = ({
    packetKey,
} : {packetKey: keyof CustomerPacketValues}) => {
    const {packet, editing, getInputProps} = useContext(PacketContext);

    const title = transformKeys[packetKey];
    const value = String(packet[packetKey]);
    return (
        <Item
            title={title}
            value={value}
            type='string'
            editComponent={
                <FormNumberInput
                    formKey={packetKey}
                    getInputProps={getInputProps}
                    key={packetKey}
                    placeholder={title}
                />
            }
            editing={editing}
        />
    );
};

type PacketMultiSelectParams = {
    packetKey: keyof CustomerPacketValues,
    data: {value: string, label: string}[]
}
export const PacketSelect = ({
    packetKey,
    data,
} : PacketMultiSelectParams) => {
    const {packet, editing, getInputProps} = useContext(PacketContext);

    const title = transformKeys[packetKey];
    const value = String(packet[packetKey]);
    return (
        <Item
            title={title}
            value={value}
            type='string'
            editComponent={
                <FormSelect
                    formKey={packetKey}
                    getInputProps={getInputProps}
                    key={packetKey}
                    placeholder={title}
                    data={data}
                />
            }
            editing={editing}
        />
    );
};

export const ProductsInUse = () => {
    const {packet: {productsInUse}, editing, getInputProps} = useContext(PacketContext);

    const title = transformKeys.productsInUse;
    const value = productsInUse.join(', ');

    return (
        <Item
            title={title}
            value={value}
            type='string'
            editComponent={
                <FormMultiSelect
                    formKey={'productsInUse'}
                    getInputProps={getInputProps}
                    key={'productsInUse'}
                    placeholder={title}
                    value={parseJson<string[]>(getInputProps('productsInUse').value)}
                    data={[
                        {value: 'calls', label: 'Calls'},
                        {value: 'playbooks', label: 'Playbooks'},
                        {value: 'boards', label: 'Boards'},
                    ]}
                />
            }
            editing={editing}
        />
    );
};
