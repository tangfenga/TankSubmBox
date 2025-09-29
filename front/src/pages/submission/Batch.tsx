import React from 'react';
import { RouteComponentProps } from 'react-router-dom';
import "./Batch.less";
import TankComponent from '../../common/component/TankComponent';
import { Button, Form, Input, message, Select, Upload } from 'antd';
import TankTitle from '../widget/TankTitle';
import TankContentCard from '../widget/TankContentCard';
import Lang from '../../common/model/global/Lang';
import Track from '../../common/model/track/Track';
import College from '../../common/model/college/College';
import { UploadOutlined } from '@ant-design/icons';
import { RcFile, UploadFile } from 'antd/es/upload/interface';
import HttpBase from '../../common/model/base/HttpBase';
import { FormInstance } from 'antd/lib/form';
import Moon from '../../common/model/global/Moon';
import Matter from '../../common/model/matter/Matter';
import UploadMatterPanel from '../matter/widget/UploadMatterPanel';

interface IProps extends RouteComponentProps {}

interface IState {
  tracks: Track[];
  colleges: College[];
  fileList: UploadFile[];
}

export default class SubmissionBatch extends TankComponent<IProps, IState> {
  trackModel = new Track(this);
  collegeModel = new College(this);
  httpBase: HttpBase = new HttpBase(this);
  formRef = React.createRef<FormInstance>();

  constructor(props: IProps) {
    super(props);
    this.state = {
      tracks: [],
      colleges: [],
      fileList: [],
    };
  }

  componentDidMount() {
    this.trackModel.httpList(() => {
      this.setState({ tracks: this.trackModel.tracks });
    });
    this.collegeModel.httpList(() => {
      this.setState({ colleges: this.collegeModel.colleges });
    });
  }

  handleSubmit = (values: any) => {
    if (this.state.fileList.length === 0) {
      message.error('请选择要上传的文件夹');
      return;
    }

    const folderData = new FormData();
    Object.keys(values).forEach((key) => {
      folderData.append(key, values[key]);
    });

    this.httpBase.httpPost(
      '/api/submission/admin/create_submission_folder',
      folderData,
      (response: any) => {
        const submissionFolder = new Matter();
        submissionFolder.assign(response.data.data);

        const puuid = submissionFolder.uuid!;
        const spaceUuid = submissionFolder.spaceUuid!;
        const ownerUuid = submissionFolder.userUuid!;

        this.uploadDirectory(puuid, spaceUuid, ownerUuid).catch(err => {
          console.error('An error occurred during directory upload:', err);
          message.error('文件夹上传过程中发生未知错误。');
        });
      },
      (error: any) => {
        this.httpBase.defaultErrorHandler(error);
      }
    );
  };

  uploadDirectory = async (puuid: string, spaceUuid: string, ownerUuid: string) => {
    const dirPathUuidMap: { [key: string]: string } = {};
    const uploadMapKey = `batch-upload-${new Date().getTime()}`;
    Moon.getSingleton().uploadMattersMap[uploadMapKey] = [];
    this.updateUI();

    let successCount = 0;
    const totalFiles = this.state.fileList.length;

    fileLoop: for (const file of this.state.fileList) {
      const originFile = file.originFileObj as File;
      if (!originFile) continue;

      const webkitRelativePath = (originFile as any).webkitRelativePath || originFile.name;
      const paths = webkitRelativePath.split('/');
      const fileName = paths.pop()!;
      
      // Skip if filename is empty (could be a directory entry)
      if (!fileName) continue;

      const pPaths = paths;

      let targetPuuid = puuid;

      if (pPaths.length > 0 && pPaths.some((p: string) => p)) {
        const pPathStr = pPaths.join('/');
        if (dirPathUuidMap[pPathStr]) {
          targetPuuid = dirPathUuidMap[pPathStr];
        } else {
          let currentPuuid = puuid;
          for (let j = 1; j <= pPaths.length; j++) {
            const midPaths = pPaths.slice(0, j);
            const midPathStr = midPaths.join('/');
            if (!dirPathUuidMap[midPathStr]) {
              const m = new Matter(this);
              m.name = midPaths[midPaths.length - 1];
              
              if (!m.name) continue;

              m.puuid = currentPuuid;
              m.spaceUuid = spaceUuid;
              m.privacy = true;

              let created = false;
              await m.httpCreateDirectory(
                () => {
                  dirPathUuidMap[midPathStr] = m.uuid!;
                  currentPuuid = m.uuid!;
                  created = true;
                },
                (msg: string) => {
                  message.error(`创建文件夹 ${m.name} 失败: ${msg}`);
                },
                ownerUuid
              );

              if (!created) {
                message.error(`无法创建文件夹 ${m.name}，中止上传该文件。`);
                continue fileLoop; 
              }
            } else {
              currentPuuid = dirPathUuidMap[midPathStr];
            }
          }
          targetPuuid = currentPuuid;
        }
      }

      const matter = new Matter(this);
      matter.file = originFile;
      matter.puuid = targetPuuid;
      matter.spaceUuid = spaceUuid;
      matter.privacy = true;
      Moon.getSingleton().uploadMattersMap[uploadMapKey].push(matter);

      matter.httpUpload(
        ownerUuid,
        () => {
          successCount++;
          const matters = Moon.getSingleton().uploadMattersMap[uploadMapKey];
          const index = matters.findIndex((m) => m === matter);
          if (index > -1) {
            matters.splice(index, 1);
          }
          if (successCount === totalFiles) {
            message.success('所有文件上传成功！');
            this.formRef.current?.resetFields();
            this.setState({ fileList: [] });
          }
        },
        (errMsg: string) => {
          message.error(`${originFile.name} 上传失败: ${errMsg}`);
        }
      );
    }
  };


  render() {
    const { tracks, colleges, fileList } = this.state;

    const uploadProps = {
      onRemove: (file: UploadFile) => {
        this.setState((prevState) => {
          const index = prevState.fileList.indexOf(file);
          const newFileList = prevState.fileList.slice();
          newFileList.splice(index, 1);
          return {
            fileList: newFileList,
          };
        });
      },
      onChange: (info: any) => {
        this.setState({ fileList: info.fileList });
      },
      beforeUpload: () => false, // Prevent auto upload
      fileList,
      directory: true,
    };

    return (
      <div className="page-submission-batch">
        <TankTitle name="为他人提交作品" />
        <TankContentCard>
          <Form
            ref={this.formRef}
            layout="vertical"
            onFinish={this.handleSubmit}
          >
            <Form.Item
              name="username"
              label="用户名"
              rules={[{ required: true, message: '请输入用户名' }]}
            >
              <Input placeholder="如果用户不存在，将使用此用户名和默认密码创建" />
            </Form.Item>
            <Form.Item
              name="realName"
              label="姓名"
              rules={[{ required: true, message: '请输入姓名' }]}
            >
              <Input />
            </Form.Item>
            <Form.Item
              name="studentId"
              label="学号/工号"
              rules={[{ required: true, message: '请输入学号/工号' }]}
            >
              <Input />
            </Form.Item>
            <Form.Item
              name="phoneNumber"
              label="手机号"
              rules={[{ required: true, message: '请输入手机号' }]}
            >
              <Input />
            </Form.Item>
            <Form.Item
              name="userType"
              label="身份"
              rules={[{ required: true, message: '请选择身份' }]}
            >
              <Select placeholder="请选择身份">
                <Select.Option value="student">学生</Select.Option>
                <Select.Option value="teacher">教师</Select.Option>
              </Select>
            </Form.Item>
            <Form.Item
              name="college"
              label="学院"
              rules={[{ required: true, message: '请选择学院' }]}
            >
              <Select placeholder="请选择学院">
                {colleges.map((c: College) => (
                  <Select.Option key={c.id} value={c.name}>
                    {c.name}
                  </Select.Option>
                ))}
              </Select>
            </Form.Item>
            <Form.Item
              name="trackId"
              label="赛道"
              rules={[{ required: true, message: '请选择赛道' }]}
            >
              <Select placeholder="请选择赛道">
                {tracks.map((track: Track) => (
                  <Select.Option key={track.id} value={track.id}>
                    {track.name}
                  </Select.Option>
                ))}
              </Select>
            </Form.Item>
            <Form.Item
              name="workName"
              label="作品名称"
              rules={[{ required: true, message: '请输入作品名称' }]}
            >
              <Input />
            </Form.Item>
            <Form.Item label="作品文件夹" required>
              <Upload {...uploadProps}>
                <Button icon={<UploadOutlined />}>选择文件夹</Button>
              </Upload>
            </Form.Item>
            <Form.Item>
              <Button
                type="primary"
                htmlType="submit"
                loading={this.httpBase.loading}
              >
                提交
              </Button>
            </Form.Item>
          </Form>
          {Object.keys(Moon.getSingleton().uploadMattersMap).map((key) =>
            Moon.getSingleton().uploadMattersMap[key].map((matter) => (
              <UploadMatterPanel key={matter.autoId} matter={matter} />
            ))
          )}
        </TankContentCard>
      </div>
    );
  }
}
