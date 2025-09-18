import React from 'react';
import { Modal } from 'antd';
import TankComponent from '../../../common/component/TankComponent';
import Lang from '../../../common/model/global/Lang';
import BottomLayout from '../BottomLayout';
import DefaultLogoPng from '../../../assets/image/logo.png';
import './AboutModal.less';
import Moon from '../../../common/model/global/Moon';

interface IProps {
  onSuccess: () => any;
  onClose: () => any;
}

interface IState {}

export default class AboutModal extends TankComponent<IProps, IState> {
  constructor(props: IProps) {
    super(props);
  }

  static open = () => {
    const modal = Modal.confirm({
      className: 'about-modal',
      title: Lang.t('layout.about'),
      width: '90vw',
      okCancel: false,
      okText: Lang.t('confirm'),
      maskClosable: true,
      content: (
        <AboutModal
          onSuccess={() => {
            modal.destroy();
          }}
          onClose={() => {
            modal.destroy();
          }}
        />
      ),
    });
  };

  changeLang() {
    BottomLayout.changeLang();
    this.updateUI();
  }

  render() {
    const { preference } = Moon.getSingleton();
    return (
      <div className="about-modal-box">
        <span className="item">
          <span dangerouslySetInnerHTML={{ __html: preference.copyright }} />
        </span>
        <span className="item">
          <span dangerouslySetInnerHTML={{ __html: preference.record }} />
        </span>
        <div className="attribution">
          <p>兰州大学萃英在线网络技术部基于以下开源项目开发：</p>
          <p>
            <a target="_blank" href="https://github.com/eyebluecn/tank">
              Tank - https://github.com/eyebluecn/tank
            </a>
          </p>
          <p>
            <a target="_blank" href="https://github.com/eyebluecn/tank-front">
              Tank Front - https://github.com/eyebluecn/tank-front
            </a>
          </p>
          <p>原项目协议：MIT License</p>
          <p>Copyright (c) 2017 eyebluecn</p>
          <p>
            Permission is hereby granted, free of charge, to any person obtaining a copy
            of this software and associated documentation files (the "Software"), to deal
            in the Software without restriction, including without limitation the rights
            to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
            copies of the Software, and to permit persons to whom the Software is
            furnished to do so, subject to the following conditions:
          </p>
          <p>
            The above copyright notice and this permission notice shall be included in all
            copies or substantial portions of the Software.
          </p>
          <p>
            THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
            IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
            FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
            AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
            LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
            OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
            SOFTWARE.
          </p>
        </div>
      </div>
    );
  }
}
