/* Autogenerated file. Do not edit manually. */
/* tslint:disable */
/* eslint-disable */
import type {
  BaseContract,
  BytesLike,
  FunctionFragment,
  Result,
  Interface,
  EventFragment,
  ContractRunner,
  ContractMethod,
  Listener,
} from "ethers";
import type {
  TypedContractEvent,
  TypedDeferredTopicFilter,
  TypedEventLog,
  TypedLogDescription,
  TypedListener,
  TypedContractMethod,
} from "../common";

export interface IECPDKSAP_MetaAddressRegistryInterface extends Interface {
  getFunction(
    nameOrSignature: "registerMetaAddress" | "resolve"
  ): FunctionFragment;

  getEvent(nameOrSignatureOrTopic: "MetaAddressRegistered"): EventFragment;

  encodeFunctionData(
    functionFragment: "registerMetaAddress",
    values: [string, BytesLike]
  ): string;
  encodeFunctionData(functionFragment: "resolve", values: [string]): string;

  decodeFunctionResult(
    functionFragment: "registerMetaAddress",
    data: BytesLike
  ): Result;
  decodeFunctionResult(functionFragment: "resolve", data: BytesLike): Result;
}

export namespace MetaAddressRegisteredEvent {
  export type InputTuple = [id: string, metaAddress: BytesLike];
  export type OutputTuple = [id: string, metaAddress: string];
  export interface OutputObject {
    id: string;
    metaAddress: string;
  }
  export type Event = TypedContractEvent<InputTuple, OutputTuple, OutputObject>;
  export type Filter = TypedDeferredTopicFilter<Event>;
  export type Log = TypedEventLog<Event>;
  export type LogDescription = TypedLogDescription<Event>;
}

export interface IECPDKSAP_MetaAddressRegistry extends BaseContract {
  connect(runner?: ContractRunner | null): IECPDKSAP_MetaAddressRegistry;
  waitForDeployment(): Promise<this>;

  interface: IECPDKSAP_MetaAddressRegistryInterface;

  queryFilter<TCEvent extends TypedContractEvent>(
    event: TCEvent,
    fromBlockOrBlockhash?: string | number | undefined,
    toBlock?: string | number | undefined
  ): Promise<Array<TypedEventLog<TCEvent>>>;
  queryFilter<TCEvent extends TypedContractEvent>(
    filter: TypedDeferredTopicFilter<TCEvent>,
    fromBlockOrBlockhash?: string | number | undefined,
    toBlock?: string | number | undefined
  ): Promise<Array<TypedEventLog<TCEvent>>>;

  on<TCEvent extends TypedContractEvent>(
    event: TCEvent,
    listener: TypedListener<TCEvent>
  ): Promise<this>;
  on<TCEvent extends TypedContractEvent>(
    filter: TypedDeferredTopicFilter<TCEvent>,
    listener: TypedListener<TCEvent>
  ): Promise<this>;

  once<TCEvent extends TypedContractEvent>(
    event: TCEvent,
    listener: TypedListener<TCEvent>
  ): Promise<this>;
  once<TCEvent extends TypedContractEvent>(
    filter: TypedDeferredTopicFilter<TCEvent>,
    listener: TypedListener<TCEvent>
  ): Promise<this>;

  listeners<TCEvent extends TypedContractEvent>(
    event: TCEvent
  ): Promise<Array<TypedListener<TCEvent>>>;
  listeners(eventName?: string): Promise<Array<Listener>>;
  removeAllListeners<TCEvent extends TypedContractEvent>(
    event?: TCEvent
  ): Promise<this>;

  registerMetaAddress: TypedContractMethod<
    [_id: string, _metaAddress: BytesLike],
    [void],
    "payable"
  >;

  resolve: TypedContractMethod<[_id: string], [string], "view">;

  getFunction<T extends ContractMethod = ContractMethod>(
    key: string | FunctionFragment
  ): T;

  getFunction(
    nameOrSignature: "registerMetaAddress"
  ): TypedContractMethod<
    [_id: string, _metaAddress: BytesLike],
    [void],
    "payable"
  >;
  getFunction(
    nameOrSignature: "resolve"
  ): TypedContractMethod<[_id: string], [string], "view">;

  getEvent(
    key: "MetaAddressRegistered"
  ): TypedContractEvent<
    MetaAddressRegisteredEvent.InputTuple,
    MetaAddressRegisteredEvent.OutputTuple,
    MetaAddressRegisteredEvent.OutputObject
  >;

  filters: {
    "MetaAddressRegistered(string,bytes)": TypedContractEvent<
      MetaAddressRegisteredEvent.InputTuple,
      MetaAddressRegisteredEvent.OutputTuple,
      MetaAddressRegisteredEvent.OutputObject
    >;
    MetaAddressRegistered: TypedContractEvent<
      MetaAddressRegisteredEvent.InputTuple,
      MetaAddressRegisteredEvent.OutputTuple,
      MetaAddressRegisteredEvent.OutputObject
    >;
  };
}