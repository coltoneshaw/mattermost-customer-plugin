import React, {createContext, useContext} from 'react';

import {GetInputProps} from '@mantine/form/lib/types';

import {CustomerPacketValues} from '@/types/customers';

import {FormTextInput} from '@/components/form/FormTextInput';

import {FormDropdown} from '@/components/form/FormDropdown';

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
    metrics: 'Metrics',
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
    metrics: false,
    metricService: '',
    samlProvider: '',
    version: '',
    mobileApp: '',
    licensedTo: '',
    productsInUse: '',
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
export const LicensedTo = () => {
    const {packet: {licensedTo}, editing, getInputProps} = useContext(PacketContext);
    const title = transformKeys.licensedTo;
    const value = licensedTo;

    return (
        <Item
            title={title}
            value={value}
            editing={editing}
            editComponent={
                <FormTextInput
                    formKey={'licensedTo'}
                    getInputProps={getInputProps}
                    key={'licensedTo'}
                    placeholder={value}
                />
            }
        />
    );
};

export const MobileApp = () => {
    const {packet: {mobileApp}, editing, getInputProps} = useContext(PacketContext);

    const title = transformKeys.mobileApp;
    const value = mobileApp;
    return (
        <Item
            title={title}
            value={value}
            editing={editing}
            editComponent={
                <FormDropdown
                    formKey={'mobileApp'}
                    getInputProps={getInputProps}
                    key={'mobileApp'}
                    placeholder={value}
                    data={[
                        {value: 'true', label: 'true'},
                        {value: 'false', label: 'false'},
                    ]}
                />
            }
        />
    );
};

export const ProductsInUse = () => {
    const {packet: {productsInUse}, editing, getInputProps} = useContext(PacketContext);

    const title = transformKeys.productsInUse;
    const value = productsInUse;
    return (
        <Item
            title={title}
            value={value}
            editComponent={
                <FormTextInput
                    formKey={'productsInUse'}
                    getInputProps={getInputProps}
                    key={'productsInUse'}
                    placeholder={value}
                />
            }
            editing={editing}
        />
    );
};

