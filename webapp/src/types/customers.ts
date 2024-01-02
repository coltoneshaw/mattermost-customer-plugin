import {AdminConfig} from '@mattermost/types/config';

// eslint-disable-next-line no-shadow
export enum LicenseType {
    Cloud = 'cloud',
    Enterprise = 'enterprise',
    Professional = 'professional',
    Free = 'free',
    Trial = 'trial',
    Nonprofit = 'nonprofit',
    Other = 'other'
}

export type Customer = {
    id: string;
    name: string;
    customerSuccessManager: string;
    accountExecutive: string;
    technicalAccountManager: string;
    salesforceId: string;
    zendeskId: string;
    type: LicenseType;
    licensedTo: string;
    siteURL: string;
    customerChannel: string;
    GDriveLink: string;
    lastUpdated: number;
}

export type CustomerPacketValues = {
    licensedTo: string;
    version: string;
    serverOS: string;
    serverArch: string;
    databaseType: string;
    databaseVersion: string;
    databaseSchemaVersion: string;
    fileDriver: string;
    activeUsers: number;
    dailyActiveUsers: number;
    monthlyActiveUsers: number;
    inactiveUserCount: number;
    licenseSupportedUsers: number;
    totalPosts: number;
    totalChannels: number;
    totalTeams: number;
}

export type CustomerPluginValues = {
    pluginID: string;
    version: string;
    isActive: boolean,
    name: string
}

export type FullCustomerInfo = Customer & {
    packet: CustomerPacketValues;
    config: AdminConfig;
    plugins: CustomerPluginValues[]
}

export type GetCustomerResult = {
    totalCount: number;
    pageCount: number;
    hasMore: boolean;
    customers: FullCustomerInfo[]
}

// eslint-disable-next-line no-shadow
export enum CustomerSortOptions {
    SortByName = 'name',
    SortByCSM = 'csm',
    SortByAE = 'ae',
    SortByTAM = 'tam',
    SortByType = 'type',
    SortBySiteURL = 'site_url',
    SortByLicensedTo = 'licensed_to',
    Default = ''
}

// eslint-disable-next-line no-shadow
export enum SortDirection {
    DirectionAsc = 'asc',
    DirectionDesc = 'desc'
}

export type CustomerFilterOptions = {
    sort: CustomerSortOptions;
    order: SortDirection;
    page: string;
    perPage: string;
    searchTerm: string;
}
