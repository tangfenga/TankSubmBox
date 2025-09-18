import React from 'react';
import TankComponent from '../../../common/component/TankComponent';
import Matter from '../../../common/model/matter/Matter';
import { Modal, InputNumber, Input, Form, Button } from 'antd';
import HttpUtil from '../../../common/util/HttpUtil';
import MessageBoxUtil from '../../../common/util/MessageBoxUtil';

interface IProps {
  matter: Matter;
  submissionId: number;
  visible: boolean;
  onCancel: () => void;
  onSuccess: () => void;
}

interface RatingItem {
  name: string;
  min: number;
  max: number;
  score: number | null;
}

interface IState {
  loading: boolean;
  ratingItems: RatingItem[];
}

export default class RatingModal extends TankComponent<IProps, IState> {
  constructor(props: IProps) {
    super(props);
    this.state = {
      loading: false,
      ratingItems: [
        { name: '创新点', min: 0, max: 20, score: null },
        { name: '完成度', min: 0, max: 20, score: null },
        { name: '美观设计', min: 0, max: 20, score: null },
      ]
    };
  }

  componentDidUpdate(prevProps: IProps) {
    if (prevProps.visible !== this.props.visible && this.props.visible) {
      // 重置表单状态
      this.setState({
        ratingItems: this.state.ratingItems.map(item => ({ ...item, score: null }))
      });
    }
  }

  // 提交评分
  submitRating = () => {
    const { ratingItems } = this.state;
    const { submissionId } = this.props;

    // 验证所有评分项
    const invalidItems = ratingItems.filter(item => 
      item.score === null || item.score < item.min || item.score > item.max
    );
    
    if (invalidItems.length > 0) {
      MessageBoxUtil.error('请填写所有评分项，并确保分数在有效范围内');
      return;
    }

    this.setState({ loading: true });

    // 为每个评分项提交评分
    let completedCount = 0;
    let hasError = false;
    
    ratingItems.forEach((item, index) => {
      const data = new FormData();
      data.set('submissionId', submissionId.toString());
      data.set('score', item.score!.toString());
      data.set('comment', item.name);
      
      HttpUtil.httpPost(
        '/api/rating/submit',
        data,
        () => {
          completedCount++;
          if (completedCount === ratingItems.length && !hasError) {
            this.setState({ loading: false });
            MessageBoxUtil.success('评分成功');
            this.props.onSuccess();
            this.props.onCancel();
          }
        },
        (msg: string) => {
          hasError = true;
          this.setState({ loading: false });
          MessageBoxUtil.error(msg || `第${index + 1}项评分失败`);
        }
      );
    });
  };

  // 处理评分变化
  handleScoreChange = (index: number, value: number | null) => {
    this.setState(prevState => ({
      ratingItems: prevState.ratingItems.map((item, i) => 
        i === index ? { ...item, score: value } : item
      )
    }));
  };

  render() {
    const { matter, visible, onCancel } = this.props;
    const { loading, ratingItems } = this.state;

    return (
      <Modal
        title={`评分 - ${matter.name}`}
        visible={visible}
        onCancel={onCancel}
        footer={[
          <Button key="cancel" onClick={onCancel}>
            取消
          </Button>,
          <Button
            key="submit"
            type="primary"
            loading={loading}
            onClick={this.submitRating}
            disabled={ratingItems.some(item => item.score === null || item.score < item.min || item.score > item.max)}
          >
            提交评分
          </Button>
        ]}
      >
        <Form layout="vertical">
          {ratingItems.map((item, index) => (
            <Form.Item 
              key={index} 
              label={`${item.name}（${item.min}-${item.max}分）`} 
              required
              validateStatus={item.score !== null && (item.score < item.min || item.score > item.max) ? 'error' : ''}
              help={item.score !== null && (item.score < item.min || item.score > item.max) ? `请输入${item.min}-${item.max}之间的分数` : ''}
            >
              <InputNumber
                min={item.min}
                max={item.max}
                value={item.score === null ? undefined : item.score}
                onChange={(value) => this.handleScoreChange(index, value)}
                style={{ width: '100%' }}
                placeholder={`请输入${item.min}-${item.max}之间的分数`}
              />
            </Form.Item>
          ))}
        </Form>
      </Modal>
    );
  }
}