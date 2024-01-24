import React from 'react';
import {NumberInput} from '@mantine/core';

import {FormTextInputParams} from './types';
import {formTextInputStyles} from './styles';

export const FormNumberInput = <T, >({
    getInputProps,
    label,
    placeholder,
    formKey,
}: FormTextInputParams<T>) => {
    const {classes} = formTextInputStyles();

    return (
        <NumberInput
            bg={'transparent'}
            classNames={classes}
            placeholder={placeholder || label}
            label={label}
            {...getInputProps(formKey)}
        />
    );
};
