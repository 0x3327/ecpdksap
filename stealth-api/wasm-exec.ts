// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

"use strict";

(() => {
	const enosys = (): Error => {
		const err = new Error("not implemented");
		(err as any).code = "ENOSYS";
		return err;
	};

	if (!globalThis.fs) {
		let outputBuf = "";
		globalThis.fs = {
			constants: { O_WRONLY: -1, O_RDWR: -1, O_CREAT: -1, O_TRUNC: -1, O_APPEND: -1, O_EXCL: -1 }, // unused
			writeSync(fd: number, buf: Uint8Array): number {
				outputBuf += decoder.decode(buf);
				const nl = outputBuf.lastIndexOf("\n");
				if (nl != -1) {
					console.log(outputBuf.substring(0, nl));
					outputBuf = outputBuf.substring(nl + 1);
				}
				return buf.length;
			},
			write(fd: number, buf: Uint8Array, offset: number, length: number, position: number | null, callback: (err: Error | null, n?: number) => void): void {
				if (offset !== 0 || length !== buf.length || position !== null) {
					callback(enosys());
					return;
				}
				const n = this.writeSync(fd, buf);
				callback(null, n);
			},
			chmod(path: string, mode: number, callback: (err: Error | null) => void): void { callback(enosys()); },
			chown(path: string, uid: number, gid: number, callback: (err: Error | null) => void): void { callback(enosys()); },
			close(fd: number, callback: (err: Error | null) => void): void { callback(enosys()); },
			fchmod(fd: number, mode: number, callback: (err: Error | null) => void): void { callback(enosys()); },
			fchown(fd: number, uid: number, gid: number, callback: (err: Error | null) => void): void { callback(enosys()); },
			fstat(fd: number, callback: (err: Error | null) => void): void { callback(enosys()); },
			fsync(fd: number, callback: (err: Error | null) => void): void { callback(null); },
			ftruncate(fd: number, length: number, callback: (err: Error | null) => void): void { callback(enosys()); },
			lchown(path: string, uid: number, gid: number, callback: (err: Error | null) => void): void { callback(enosys()); },
			link(path: string, link: string, callback: (err: Error | null) => void): void { callback(enosys()); },
			lstat(path: string, callback: (err: Error | null) => void): void { callback(enosys()); },
			mkdir(path: string, perm: number, callback: (err: Error | null) => void): void { callback(enosys()); },
			open(path: string, flags: number, mode: number, callback: (err: Error | null) => void): void { callback(enosys()); },
			read(fd: number, buffer: Uint8Array, offset: number, length: number, position: number | null, callback: (err: Error | null, bytesRead?: number) => void): void { callback(enosys()); },
			readdir(path: string, callback: (err: Error | null) => void): void { callback(enosys()); },
			readlink(path: string, callback: (err: Error | null) => void): void { callback(enosys()); },
			rename(from: string, to: string, callback: (err: Error | null) => void): void { callback(enosys()); },
			rmdir(path: string, callback: (err: Error | null) => void): void { callback(enosys()); },
			stat(path: string, callback: (err: Error | null) => void): void { callback(enosys()); },
			symlink(path: string, link: string, callback: (err: Error | null) => void): void { callback(enosys()); },
			truncate(path: string, length: number, callback: (err: Error | null) => void): void { callback(enosys()); },
			unlink(path: string, callback: (err: Error | null) => void): void { callback(enosys()); },
			utimes(path: string, atime: number, mtime: number, callback: (err: Error | null) => void): void { callback(enosys()); },
		};
	}

	if (!globalThis.process) {
		globalThis.process = {
			getuid(): number { return -1; },
			getgid(): number { return -1; },
			geteuid(): number { return -1; },
			getegid(): number { return -1; },
			getgroups(): number[] { throw enosys(); },
			pid: -1,
			ppid: -1,
			umask(): number { throw enosys(); },
			cwd(): string { throw enosys(); },
			chdir(path: string): void { throw enosys(); },
		};
	}

	if (!globalThis.crypto) {
		throw new Error("globalThis.crypto is not available, polyfill required (crypto.getRandomValues only)");
	}

	if (!globalThis.performance) {
		throw new Error("globalThis.performance is not available, polyfill required (performance.now only)");
	}

	if (!globalThis.TextEncoder) {
		throw new Error("globalThis.TextEncoder is not available, polyfill required");
	}

	if (!globalThis.TextDecoder) {
		throw new Error("globalThis.TextDecoder is not available, polyfill required");
	}

	const encoder = new TextEncoder();
	const decoder = new TextDecoder();

	class Go {
		argv: string[];
		env: { [key: string]: string };
		exit: (code: number) => void;
		_exitPromise: Promise<void>;
		_resolveExitPromise!: () => void;
		_pendingEvent: any;
		_scheduledTimeouts: Map<number, NodeJS.Timeout>;
		_nextCallbackTimeoutID: number;
		mem!: DataView;
		_inst!: WebAssembly.Instance;
		_values!: any[];
		_goRefCounts!: number[];
		_ids!: Map<any, number>;
		_idPool!: number[];
		importObject: any;
		exited: boolean = false;

		constructor() {
			this.argv = ["js"];
			this.env = {};
			this.exit = (code: number) => {
				if (code !== 0) {
					console.warn("exit code:", code);
				}
			};
			this._exitPromise = new Promise((resolve) => {
				this._resolveExitPromise = resolve;
			});
			this._pendingEvent = null;
			this._scheduledTimeouts = new Map();
			this._nextCallbackTimeoutID = 1;

			const setInt64 = (addr: number, v: number) => {
				this.mem.setUint32(addr + 0, v, true);
				this.mem.setUint32(addr + 4, Math.floor(v / 4294967296), true);
			};

			const setInt32 = (addr: number, v: number) => {
				this.mem.setUint32(addr + 0, v, true);
			};

			const getInt64 = (addr: number): number => {
				const low = this.mem.getUint32(addr + 0, true);
				const high = this.mem.getInt32(addr + 4, true);
				return low + high * 4294967296;
			};

			const loadValue = (addr: number): any => {
				const f = this.mem.getFloat64(addr, true);
				if (f === 0) {
					return undefined;
				}
				if (!isNaN(f)) {
					return f;
				}

				const id = this.mem.getUint32(addr, true);
				return this._values[id];
			};

			const storeValue = (addr: number, v: any) => {
				const nanHead = 0x7FF80000;

				if (typeof v === "number" && v !== 0) {
					if (isNaN(v)) {
						this.mem.setUint32(addr + 4, nanHead, true);
						this.mem.setUint32(addr, 0, true);
						return;
					}
					this.mem.setFloat64(addr, v, true);
					return;
				}

				if (v === undefined) {
					this.mem.setFloat64(addr, 0, true);
					return;
				}

				let id = this._ids.get(v);
				if (id === undefined) {
					id = this._idPool.pop();
					if (id === undefined) {
						id = this._values.length;
					}
					this._values[id] = v;
					this._goRefCounts[id] = 0;
					this._ids.set(v, id);
				}
				this._goRefCounts[id]++;
				let typeFlag = 0;
				switch (typeof v) {
					case "object":
                        if (v !== null) {
							typeFlag = 1;
						}
						break;
					case "string":
						typeFlag = 2;
						break;
					case "symbol":
						typeFlag = 3;
						break;
					case "function":
						typeFlag = 4;
						break;
				}
				this.mem.setUint32(addr + 4, nanHead | typeFlag, true);
				this.mem.setUint32(addr, id, true);
			};

			const loadSlice = (addr: number): Uint8Array => {
				const array = getInt64(addr + 0);
				const len = getInt64(addr + 8);
				return new Uint8Array(this.mem.buffer, array, len);
			};

			const loadSliceOfValues = (addr: number): any[] => {
				const array = getInt64(addr + 0);
				const len = getInt64(addr + 8);
				const a = new Array(len);
				for (let i = 0; i < len; i++) {
					a[i] = loadValue(array + i * 8);
				}
				return a;
			};

			const loadString = (addr: number): string => {
				const saddr = getInt64(addr + 0);
				const len = getInt64(addr + 8);
				return decoder.decode(new DataView(this.mem.buffer, saddr, len));
			};

			const timeOrigin = Date.now() - performance.now();

			this.importObject = {
				wasi_snapshot_preview1: {
					fd_close: (fd: number): number => {
						return 0;
					},
					fd_fdstat_get: (fd: number, buf: number): number => {
						return 0;
					},
					fd_seek: (fd: number, offset_low: number, offset_high: number, whence: number, newOffset: number): number => {
						return 0;
					},
					fd_write: (fd: number, iovs_ptr: number, iovs_len: number, nwritten_ptr: number): number => {
						let nwritten = 0;
						for (let i = 0; i < iovs_len; i++) {
							const iov_ptr = iovs_ptr + i * 8;
							const ptr = getInt64(iov_ptr + 0);
							const len = getInt64(iov_ptr + 8);
							nwritten += len;
							if (fd === 1 || fd === 2) {
								const s = decoder.decode(new DataView(this.mem.buffer, ptr, len));
								console.log(s);
							}
						}
						setInt64(nwritten_ptr, nwritten);
						return 0;
					},
				},
				env: {
					"syscall/js.finalizeRef": (v_addr: number): void => {
						const id = this.mem.getUint32(v_addr, true);
						this._goRefCounts[id]--;
						if (this._goRefCounts[id] === 0) {
							const v = this._values[id];
							this._values[id] = null;
							this._ids.delete(v);
							this._idPool.push(id);
						}
					},
					"syscall/js.stringVal": (ret_ptr: number, value_ptr: number, value_len: number): void => {
						const s = decoder.decode(new DataView(this.mem.buffer, value_ptr, value_len));
						storeValue(ret_ptr, s);
					},
					"syscall/js.valueGet": (v_addr: number, p_ptr: number, p_len: number, ret_ptr: number): void => {
						const prop = decoder.decode(new DataView(this.mem.buffer, p_ptr, p_len));
						const v = loadValue(v_addr);
						storeValue(ret_ptr, Reflect.get(v, prop));
					},
					"syscall/js.valueSet": (v_addr: number, p_ptr: number, p_len: number, x_addr: number): void => {
						const prop = decoder.decode(new DataView(this.mem.buffer, p_ptr, p_len));
						const v = loadValue(v_addr);
						const x = loadValue(x_addr);
						Reflect.set(v, prop, x);
					},
					"syscall/js.valueIndex": (v_addr: number, i: number, ret_ptr: number): void => {
						storeValue(ret_ptr, Reflect.get(loadValue(v_addr), i));
					},
					"syscall/js.valueSetIndex": (v_addr: number, i: number, x_addr: number): void => {
						Reflect.set(loadValue(v_addr), i, loadValue(x_addr));
					},
					"syscall/js.valueCall": (
						v_addr: number,
						m_ptr: number,
						m_len: number,
						args_ptr: number,
						args_len: number,
						ret_ptr: number,
					): void => {
						const method = decoder.decode(new DataView(this.mem.buffer, m_ptr, m_len));
						const v = loadValue(v_addr);
						const args = loadSliceOfValues(args_ptr);
						try {
							const result = Reflect.apply(Reflect.get(v, method), v, args);
							storeValue(ret_ptr, result);
							this.mem.setUint8(ret_ptr + 8, 1);
						} catch (err) {
							storeValue(ret_ptr, err);
							this.mem.setUint8(ret_ptr + 8, 0);
						}
					},
					"syscall/js.valueInvoke": (v_addr: number, args_ptr: number, args_len: number, ret_ptr: number): void => {
						const v = loadValue(v_addr);
						const args = loadSliceOfValues(args_ptr);
						try {
							const result = Reflect.apply(v, undefined, args);
							storeValue(ret_ptr, result);
							this.mem.setUint8(ret_ptr + 8, 1);
						} catch (err) {
							storeValue(ret_ptr, err);
							this.mem.setUint8(ret_ptr + 8, 0);
						}
					},
					"syscall/js.valueNew": (v_addr: number, args_ptr: number, args_len: number, ret_ptr: number): void => {
						const v = loadValue(v_addr);
						const args = loadSliceOfValues(args_ptr);
						try {
							const result = Reflect.construct(v, args);
							storeValue(ret_ptr, result);
							this.mem.setUint8(ret_ptr + 8, 1);
						} catch (err) {
							storeValue(ret_ptr, err);
							this.mem.setUint8(ret_ptr + 8, 0);
						}
					},
					"syscall/js.valueLength": (v_addr: number): number => {
						return Reflect.get(loadValue(v_addr), "length");
					},
					"syscall/js.valuePrepareString": (v_addr: number, ret_ptr: number): void => {
						const s = String(loadValue(v_addr));
						const str = encoder.encode(s);
						storeValue(ret_ptr, str);
						setInt64(ret_ptr + 8, str.length);
					},
					"syscall/js.valueLoadString": (v_addr: number, slice_ptr: number): void => {
						const str = loadValue(v_addr);
						const array = getInt64(slice_ptr + 0);
						const len = getInt64(slice_ptr + 8);
						(new Uint8Array(this.mem.buffer, array, len)).set(str);
					},
					"syscall/js.valueInstanceOf": (v_addr: number, t_addr: number): number => {
						return loadValue(v_addr) instanceof loadValue(t_addr) ? 1 : 0;
					},
					"syscall/js.copyBytesToGo": (dst_addr: number, src_addr: number, src_len: number): number => {
						const dst = new Uint8Array(this.mem.buffer, getInt64(dst_addr), src_len);
						const src = new Uint8Array(this.mem.buffer, src_addr, src_len);
						dst.set(src);
						return src_len;
					},
					"syscall/js.copyBytesToJS": (dst_addr: number, src_addr: number, dst_len: number): number => {
						const dst = new Uint8Array(this.mem.buffer, dst_addr, dst_len);
						const src = new Uint8Array(this.mem.buffer, getInt64(src_addr), dst_len);
						dst.set(src);
						return dst_len;
					},
				},
			};
		}

		async run(instance: WebAssembly.Instance): Promise<void> {
			this._inst = instance;
			this.mem = new DataView(this._inst.exports.mem.buffer);
			(this._inst.exports.run as Function)();
			if (this.exited) {
				return;
			}
			await this._exitPromise;
		}

		// async importObject(): Promise<WebAssembly.Imports> {
		// 	return this.importObject;
		// }
	}

	globalThis.Go = Go;
})();

